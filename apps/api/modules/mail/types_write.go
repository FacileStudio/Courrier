package mail

// DraftRequest is the body of the draft endpoints. It is a SendRequest plus,
// on update, the ids of stored attachments to drop; the attachments not listed
// are kept.
type DraftRequest struct {
	SendRequest
	RemoveAttachments []int64 `json:"remove_attachments"`
}

// MoveEmailsRequest is the body of POST /emails/move. Name the destination
// with folder_id, or with folder_type for a system folder (inbox, sent,
// drafts, trash, junk, archive).
type MoveEmailsRequest struct {
	EmailIDs   []int64 `json:"email_ids"`
	FolderID   *int64  `json:"folder_id"`
	FolderType string  `json:"folder_type"`
}

// FolderNameRequest is the body of the folder create and rename endpoints.
type FolderNameRequest struct {
	Name string `json:"name"`
}

// DeleteFolderResponse confirms a folder deletion.
type DeleteFolderResponse struct {
	Deleted bool `json:"deleted"`
}
