package mail

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// CreateFolder creates a top-level mailbox on the IMAP server and caches it.
// The name must not collide with a system folder, because the next sync would
// classify it as one.
func (s *Service) CreateFolder(ctx context.Context, userID, accountID int64, name string) (schemas.Folder, error) {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return schemas.Folder{}, err
	}
	name = strings.TrimSpace(name)
	if err := validateFolderName(name, "/"); err != nil {
		return schemas.Folder{}, err
	}
	if s.pathTaken(ctx, accountID, name, 0) {
		return schemas.Folder{}, errors.Conflict("a folder with that name already exists")
	}

	client, err := s.connect(account)
	if err != nil {
		return schemas.Folder{}, errors.Failed(err.Error())
	}
	defer closeIMAP(client)
	if err := validateFolderName(name, folderDelimiter(client)); err != nil {
		return schemas.Folder{}, err
	}
	if err := client.Create(name, nil).Wait(); err != nil {
		return schemas.Folder{}, errors.Failed("the server refused to create the folder: " + err.Error())
	}

	folder := schemas.Folder{AccountID: accountID, SpaceID: account.SpaceID, Path: name, Name: name, Type: schemas.FolderTypeCustom}
	if err := s.orm.WithContext(ctx).Create(&folder).Error; err != nil {
		return schemas.Folder{}, errors.Internal("failed to save folder", err)
	}
	return folder, nil
}

// RenameFolder renames a custom mailbox. IMAP renames the children with it, so
// the cached paths of every descendant are rewritten to match.
func (s *Service) RenameFolder(ctx context.Context, userID, accountID, folderID int64, name string) (schemas.Folder, error) {
	account, folder, err := s.customFolder(ctx, userID, accountID, folderID, "renamed")
	if err != nil {
		return folder, err
	}
	name = strings.TrimSpace(name)
	if err := validateFolderName(name, "/"); err != nil {
		return folder, err
	}
	client, err := s.connect(account)
	if err != nil {
		return folder, errors.Failed(err.Error())
	}
	defer closeIMAP(client)

	delim := folderDelimiter(client)
	if err := validateFolderName(name, delim); err != nil {
		return folder, err
	}
	newPath := siblingPath(folder.Path, name, delim)
	if s.pathTaken(ctx, accountID, newPath, folder.ID) {
		return folder, errors.Conflict("a folder with that name already exists")
	}
	if err := client.Rename(folder.Path, newPath, nil).Wait(); err != nil {
		return folder, errors.Failed("the server refused to rename the folder: " + err.Error())
	}
	return s.repath(ctx, folder, newPath, name, delim)
}

func (s *Service) repath(ctx context.Context, folder schemas.Folder, newPath, name, delim string) (schemas.Folder, error) {
	err := s.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var children []schemas.Folder
		tx.Where("account_id = ? AND path LIKE ?", folder.AccountID, escapeLikePattern(folder.Path+delim)+"%").Find(&children)
		for _, child := range children {
			path := newPath + child.Path[len(folder.Path):]
			if err := tx.Model(&schemas.Folder{}).Where("id = ?", child.ID).Update("path", path).Error; err != nil {
				return err
			}
		}
		return tx.Model(&folder).Updates(map[string]any{"path": newPath, "name": name}).Error
	})
	if err != nil {
		return folder, errors.Internal("failed to save folder", err)
	}
	folder.Path, folder.Name = newPath, name
	return folder, nil
}

// DeleteFolder deletes a custom mailbox from the server and the cache. IMAP
// deletes the messages inside with it, so a folder that holds messages is
// refused unless force is set.
func (s *Service) DeleteFolder(ctx context.Context, userID, accountID, folderID int64, force bool) error {
	account, folder, err := s.customFolder(ctx, userID, accountID, folderID, "deleted")
	if err != nil {
		return err
	}
	var cached int64
	s.orm.WithContext(ctx).Model(&schemas.Email{}).Where("folder_id = ?", folder.ID).Count(&cached)
	if !force && (cached > 0 || folder.TotalCount > 0) {
		return errors.Conflict("the folder still holds messages: move them first, or pass force=true to delete them with it")
	}
	client, err := s.connect(account)
	if err != nil {
		return errors.Failed(err.Error())
	}
	defer closeIMAP(client)
	if err := client.Delete(folder.Path).Wait(); err != nil {
		return errors.Failed("the server refused to delete the folder: " + err.Error())
	}
	return s.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ids []int64
		tx.Model(&schemas.Email{}).Where("folder_id = ?", folder.ID).Pluck("id", &ids)
		if len(ids) > 0 {
			tx.Where("email_id IN ?", ids).Delete(&schemas.Attachment{})
			tx.Where("folder_id = ?", folder.ID).Delete(&schemas.Email{})
		}
		return tx.Delete(&folder).Error
	})
}

func (s *Service) customFolder(ctx context.Context, userID, accountID, folderID int64, verb string) (schemas.Account, schemas.Folder, error) {
	account, err := s.getAccountDecrypted(ctx, userID, accountID)
	if err != nil {
		return account, schemas.Folder{}, err
	}
	var folder schemas.Folder
	if err := s.orm.WithContext(ctx).Where("id = ? AND account_id = ?", folderID, accountID).First(&folder).Error; err != nil {
		return account, folder, errors.NotFound("folder not found")
	}
	if folder.Type != schemas.FolderTypeCustom {
		return account, folder, errors.Invalid("system folders cannot be " + verb)
	}
	return account, folder, nil
}
