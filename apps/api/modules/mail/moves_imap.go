package mail

import (
	"fmt"
	"sort"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func hasServerCopies(group []schemas.Email) bool {
	for _, e := range group {
		if e.IMAPUID != 0 {
			return true
		}
	}
	return false
}

func uidsOf(group []schemas.Email) []imap.UID {
	var uids []imap.UID
	for _, e := range group {
		if e.IMAPUID != 0 {
			uids = append(uids, imap.UID(e.IMAPUID))
		}
	}
	sort.Slice(uids, func(i, j int) bool { return uids[i] < uids[j] })
	return uids
}

// moveUIDs moves messages between two mailboxes in one command and returns the
// uid each one received in the destination. The map is empty when the server
// reports no UIDPLUS data, in which case callers store uid 0 rather than a
// stale uid that names another message or none.
func moveUIDs(client *imapclient.Client, src string, uids []imap.UID, dest string) (map[uint32]uint32, error) {
	if _, err := client.Select(src, nil).Wait(); err != nil {
		return nil, fmt.Errorf("SELECT %q failed: %w", src, err)
	}
	data, err := client.Move(imap.UIDSetNum(uids...), dest).Wait()
	if err != nil {
		return nil, fmt.Errorf("MOVE to %q failed: %w", dest, err)
	}
	return pairUIDs(data), nil
}

func pairUIDs(data *imapclient.MoveData) map[uint32]uint32 {
	out := map[uint32]uint32{}
	from, ok1 := data.SourceUIDs.(imap.UIDSet)
	to, ok2 := data.DestUIDs.(imap.UIDSet)
	if !ok1 || !ok2 {
		return out
	}
	src, _ := from.Nums()
	dst, _ := to.Nums()
	if len(src) != len(dst) {
		return out
	}
	sort.Slice(src, func(i, j int) bool { return src[i] < src[j] })
	sort.Slice(dst, func(i, j int) bool { return dst[i] < dst[j] })
	for i := range src {
		out[uint32(src[i])] = uint32(dst[i])
	}
	return out
}

// expungeFromServer permanently removes the server copies of emails already in
// the folder, so a later sync does not bring them back. Rows with no uid have
// no server copy.
func (s *Service) expungeFromServer(account schemas.Account, folder schemas.Folder, group []schemas.Email) error {
	uids := uidsOf(group)
	if len(uids) == 0 {
		return nil
	}
	client, err := s.connect(account)
	if err != nil {
		return err
	}
	defer closeIMAP(client)
	return expungeUIDs(client, folder.Path, uids)
}
