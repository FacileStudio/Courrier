package spaces

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/FacileStudio/porte/directory"

	"gorm.io/gorm"
)

// importBackoff is the first wait after a refused import. It doubles up to
// maxImportBackoff, so a Registre that is not deployed yet costs one log line
// every few minutes rather than a crash or a flood.
const (
	importBackoff    = 30 * time.Second
	maxImportBackoff = 5 * time.Minute
)

// Directory keeps Courrier's shared spaces equal to Registre's. It hands
// Registre the spaces Courrier held before Registre owned them, once, and then
// applies every snapshot Registre serves. Mail accounts attached to a space
// are never part of it: only the space and its membership come from Registre.
type Directory struct {
	orm    *gorm.DB
	issuer string
	logger *slog.Logger
	syncer *directory.Syncer
}

// NewDirectory builds the Directory for the OIDC client Courrier signs users
// in with. The client's issuer is also the provider key of the porte identity
// rows members are matched on.
func NewDirectory(orm *gorm.DB, client directory.Client, logger *slog.Logger) *Directory {
	d := &Directory{orm: orm, issuer: client.Issuer, logger: logger}
	d.syncer = &directory.Syncer{Client: client, Apply: d.Apply, Logger: logger}
	return d
}

// WithEntitlements chains apply into every snapshot the Syncer pulls, ahead
// of the space sync, so the plan values Registre serves beside the spaces land
// with the same pull. Both run on every snapshot and both errors are reported,
// so a failing space sync never freezes the plan values or the reverse.
func (d *Directory) WithEntitlements(apply func(context.Context, directory.Snapshot) error) *Directory {
	spaces := d.syncer.Apply
	d.syncer.Apply = func(ctx context.Context, snapshot directory.Snapshot) error {
		return errors.Join(apply(ctx, snapshot), spaces(ctx, snapshot))
	}
	return d
}

// Handler is Registre's nudge endpoint, mounted at directory.SyncPath. It
// takes no credentials: a nudge only makes the next authenticated pull sooner.
func (d *Directory) Handler() http.Handler {
	return d.syncer.Handler()
}

// Run imports Courrier's spaces, retrying until Registre accepts, and only
// then syncs until ctx ends. Syncing first would remove every member Registre
// has not heard of yet, which before the import is everyone.
func (d *Directory) Run(ctx context.Context) {
	importThenSync(ctx, d.syncer, d.imports, importBackoff)
}

// importThenSync is Run with the import source and the first backoff passed
// in, so the ordering can be tested without a database.
func importThenSync(ctx context.Context, syncer *directory.Syncer,
	load func(context.Context) ([]directory.ImportSpace, error), backoff time.Duration) {
	for delay := backoff; ; delay = min(2*delay, maxImportBackoff) {
		err := importOnce(ctx, syncer.Client, load)
		if err == nil {
			break
		}
		syncer.Logger.Warn("registre: space import failed, retrying",
			slog.Any("error", err), slog.Duration("retry_in", delay))
		if !wait(ctx, delay) {
			return
		}
	}
	syncer.Run(ctx)
}

// importOnce reads the local spaces afresh and sends them, so a retry carries
// whatever changed while Registre was unreachable.
func importOnce(ctx context.Context, client directory.Client,
	load func(context.Context) ([]directory.ImportSpace, error)) error {
	spaces, err := load(ctx)
	if err != nil {
		return err
	}
	_, err = client.Import(ctx, spaces)
	return err
}

// wait sleeps for d or until ctx ends, reporting whether it slept fully.
func wait(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
