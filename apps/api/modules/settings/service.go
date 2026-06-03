package settings

import (
	"context"
	stderrors "errors"
	"strings"

	"api/internal/errors"
	"api/schemas"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const appSettingID = 1

type Service struct {
	orm        *gorm.DB
	controller *Controller
}

func NewService(orm *gorm.DB) *Service {
	service := &Service{orm: orm}
	service.controller = newController(service)
	return service
}

func (service *Service) getSettings(ctx context.Context) (*Settings, error) {
	var record schemas.AppSetting
	err := service.orm.WithContext(ctx).Where("id = ?", appSettingID).First(&record).Error
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return &Settings{}, nil
	}
	if err != nil {
		return nil, errors.Internal("failed to get settings", err)
	}
	return &Settings{
		NookPoolURL:     record.NookPoolURL,
		NookPoolSecret:  record.NookPoolSecret,
		NookPoolEnabled: record.NookPoolEnabled,
	}, nil
}

func (service *Service) updateSettings(ctx context.Context, req *UpdateRequest) (*Settings, error) {
	record := schemas.AppSetting{
		ID:              appSettingID,
		NookPoolURL:     strings.TrimSpace(req.NookPoolURL),
		NookPoolSecret:  strings.TrimSpace(req.NookPoolSecret),
		NookPoolEnabled: req.NookPoolEnabled,
	}
	if err := service.orm.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"nook_pool_url", "nook_pool_secret", "nook_pool_enabled"}),
	}).Create(&record).Error; err != nil {
		return nil, errors.Internal("failed to update settings", err)
	}
	return &Settings{
		NookPoolURL:     record.NookPoolURL,
		NookPoolSecret:  record.NookPoolSecret,
		NookPoolEnabled: record.NookPoolEnabled,
	}, nil
}
