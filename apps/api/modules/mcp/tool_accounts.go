package mcp

import (
	"context"

	mcpkit "github.com/FacileStudio/tronc/mcp"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// listAccounts is courrier_list_accounts.
func (t *toolset) listAccounts(ctx context.Context, req *sdk.CallToolRequest, in accountsInput) (*sdk.CallToolResult, accountsOutput, error) {
	who, err := callerOf(req)
	if err != nil {
		return nil, accountsOutput{}, err
	}
	full, err := mcpkit.Detailed(in.ResponseFormat)
	if err != nil {
		return nil, accountsOutput{}, err
	}
	accounts, err := t.mail.VisibleAccounts(ctx, who.userID)
	if err != nil {
		return nil, accountsOutput{}, mcpkit.Failure(err, "")
	}
	out := accountsOutput{Accounts: make([]accountSummary, 0, len(accounts))}
	for _, account := range accounts {
		out.Accounts = append(out.Accounts, accountSummary{
			AccountID: account.ID, Email: account.Email, Name: account.Name,
			Shared: account.SpaceID != nil && *account.SpaceID != "", Default: account.IsDefault,
		})
	}
	if full {
		err = t.addFolders(ctx, accounts, out.Accounts)
	}
	return nil, out, err
}

// addFolders fills each summary's folders, in the order Courrier lists them:
// system folders by type, then the person's own by name.
func (t *toolset) addFolders(ctx context.Context, accounts []schemas.Account, summaries []accountSummary) error {
	ids := make([]int64, len(accounts))
	for i, account := range accounts {
		ids[i] = account.ID
	}
	folders, err := t.mail.FoldersByID(ctx, ids)
	if err != nil {
		return mcpkit.Failure(err, "")
	}
	byAccount := map[int64][]folderSummary{}
	for _, folder := range sortedFolders(folders) {
		byAccount[folder.AccountID] = append(byAccount[folder.AccountID], folderSummary{
			FolderID: folder.ID, Type: folder.Type, Name: folder.Name, Unread: folder.UnreadCount, Total: folder.TotalCount,
		})
	}
	for i := range summaries {
		summaries[i].Folders = byAccount[summaries[i].AccountID]
	}
	return nil
}
