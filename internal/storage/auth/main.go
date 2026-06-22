package auth

import (
	"context"
	"fmt"
	"log/slog"

	"io"

	"time"

	"go.uber.org/fx"

	"github.com/dgraph-io/badger/v4"
	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/auth/transactions"
	"link-society.com/flowg/internal/storage/changefeed"
	"link-society.com/flowg/internal/storage/schema"
	"link-society.com/flowg/internal/utils/fxproviders"
	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/kvstore"
)

type Storage interface {
	storage.Streamable

	ListRoles(ctx context.Context) ([]models.Role, error)
	FetchRole(ctx context.Context, name string) (*models.Role, error)
	SaveRole(ctx context.Context, role models.Role) error
	DeleteRole(ctx context.Context, name string) error

	ListUsers(ctx context.Context) ([]models.User, error)
	FetchUser(ctx context.Context, name string) (*models.User, error)
	ListUserScopes(ctx context.Context, name string) ([]models.Scope, error)
	SaveUser(ctx context.Context, user models.User, password string) error
	PatchUserRoles(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, name string) error

	VerifyUserPassword(ctx context.Context, name, password string) (bool, error)
	VerifyUserPermission(ctx context.Context, username string, scope models.Scope) (bool, error)

	CreateToken(ctx context.Context, username string) (string, string, error)
	VerifyToken(ctx context.Context, token string) (*models.User, error)
	ListTokens(ctx context.Context, username string) ([]string, error)
	DeleteToken(ctx context.Context, username string, tokenUUID string) error
}

type Options struct {
	Directory string
	InMemory  bool
	ReadOnly  bool

	GCInterval           time.Duration
	TombstoneGracePeriod time.Duration
}

const (
	roleKind  = "role"
	userKind  = "user"
	tokenKind = "token"
)

type storageImpl struct {
	kvStore  kvstore.Storage
	clock    *hlc.Clock
	notifier changefeed.Notifier
}

type deps struct {
	fx.In

	S        kvstore.Storage `name:"storage.auth"`
	Clock    *hlc.Clock
	Notifier changefeed.Notifier
}

var _ Storage = (*storageImpl)(nil)

func DefaultOptions() Options {
	return Options{
		Directory: "",
		InMemory:  false,
		ReadOnly:  false,

		GCInterval:           time.Hour,
		TombstoneGracePeriod: 24 * time.Hour,
	}
}

func NewStorage(opts Options) fx.Option {
	kvOpts := kvstore.DefaultOptions()
	kvOpts.LogChannel = "storage.auth"
	kvOpts.Directory = opts.Directory
	kvOpts.InMemory = opts.InMemory
	kvOpts.ReadOnly = opts.ReadOnly

	return fx.Module(
		"storage.auth",
		kvstore.NewStorage(kvOpts),
		fx.Provide(func(lc fx.Lifecycle, d deps) Storage {
			storage := &storageImpl{kvStore: d.S, clock: d.Clock, notifier: d.Notifier}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := schema.Migrate(ctx, storage.kvStore, storage.clock, nil); err != nil {
						return fmt.Errorf("failed to migrate schema: %w", err)
					}

					if err := migrateAlertScopes(ctx, storage.kvStore, storage.clock); err != nil {
						return fmt.Errorf("failed to migrate alerts: %w", err)
					}

					return nil
				},
			})

			return storage
		}),
		fxproviders.ProvideActor[*gcActor](func(d deps) *gcActor {
			return &gcActor{
				Actor: actor.New(&gcWorker{
					kvStore:    d.S,
					grace:      opts.TombstoneGracePeriod,
					gcInterval: opts.GCInterval,
				}),
			}
		}),
	)
}

func (s *storageImpl) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	return s.kvStore.Backup(ctx, w, since)
}

func (s *storageImpl) DropAll(ctx context.Context) error {
	return s.kvStore.DropAll(ctx)
}

func (s *storageImpl) LatestVersion(ctx context.Context) (uint64, error) {
	return s.kvStore.LatestVersion(ctx)
}

func (s *storageImpl) Load(ctx context.Context, r io.Reader) error {
	return s.kvStore.Restore(ctx, r)
}

func (s *storageImpl) Merge(ctx context.Context, r io.Reader) error {
	var applied []schema.AppliedRecord
	if err := s.kvStore.Merge(ctx, r, schema.MergeEnveloped(&applied)); err != nil {
		return err
	}

	if len(applied) > 0 {
		s.emitResync(ctx)
	}
	return nil
}

func (s *storageImpl) ApplyReplicated(ctx context.Context, records []changefeed.Record) error {
	var applied bool
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			for _, record := range records {
				_, ok, err := schema.ApplyRecord(txn, record.Key, record.Value)
				if err != nil {
					return err
				}
				if ok {
					applied = true
				}
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	if applied {
		s.emitResync(ctx)
	}
	return nil
}

func (s *storageImpl) ListRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			roles, err = transactions.ListRoles(txn)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *storageImpl) FetchRole(ctx context.Context, name string) (*models.Role, error) {
	var role *models.Role

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			role, err = transactions.FetchRole(txn, name)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *storageImpl) SaveRole(ctx context.Context, role models.Role) error {
	ts := s.clock.Now()
	var records []changefeed.Record
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.SaveRole(txn, role, ts, &records)
		},
	)
	if err != nil {
		return err
	}

	s.emitChange(ctx, roleKind, role.Name, changefeed.OpWrite, ts.NodeID, records)
	return nil
}

func (s *storageImpl) DeleteRole(ctx context.Context, name string) error {
	ts := s.clock.Now()
	var records []changefeed.Record
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteRole(txn, name, ts, &records)
		},
	)
	if err != nil {
		return err
	}

	s.emitChange(ctx, roleKind, name, changefeed.OpDelete, ts.NodeID, records)
	return nil
}

func (s *storageImpl) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			users, err = transactions.ListUsers(txn)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storageImpl) FetchUser(ctx context.Context, name string) (*models.User, error) {
	var user *models.User

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			user, err = transactions.FetchUser(txn, name)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storageImpl) ListUserScopes(ctx context.Context, name string) ([]models.Scope, error) {
	var scopes []models.Scope

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			scopes, err = transactions.ListUserScopes(txn, name)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return scopes, nil
}

func (s *storageImpl) SaveUser(ctx context.Context, user models.User, password string) error {
	ts := s.clock.Now()
	var records []changefeed.Record
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.SaveUser(txn, user, password, ts, &records)
		},
	)
	if err != nil {
		return err
	}

	s.emitChange(ctx, userKind, user.Name, changefeed.OpWrite, ts.NodeID, records)
	return nil
}

func (s *storageImpl) PatchUserRoles(ctx context.Context, user models.User) error {
	ts := s.clock.Now()
	var records []changefeed.Record
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.PatchUserRoles(txn, user, ts, &records)
		},
	)
	if err != nil {
		return err
	}

	s.emitChange(ctx, userKind, user.Name, changefeed.OpWrite, ts.NodeID, records)
	return nil
}

func (s *storageImpl) DeleteUser(ctx context.Context, name string) error {
	ts := s.clock.Now()
	var records []changefeed.Record
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteUser(txn, name, ts, &records)
		},
	)
	if err != nil {
		return err
	}

	s.emitChange(ctx, userKind, name, changefeed.OpDelete, ts.NodeID, records)
	return nil
}

func (s *storageImpl) VerifyUserPassword(ctx context.Context, name, password string) (bool, error) {
	var verified bool

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			verified, err = transactions.VerifyUserPassword(txn, name, password)
			return err
		},
	)
	if err != nil {
		return false, err
	}

	return verified, nil
}

func (s *storageImpl) VerifyUserPermission(ctx context.Context, username string, scope models.Scope) (bool, error) {
	var authorized bool

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			authorized, err = transactions.VerifyUserPermission(txn, username, scope)
			return err
		},
	)
	if err != nil {
		return false, err
	}

	return authorized, nil
}

func (s *storageImpl) CreateToken(ctx context.Context, username string) (string, string, error) {
	var token, tokenUuid string

	ts := s.clock.Now()
	var records []changefeed.Record
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			token, tokenUuid, err = transactions.CreateToken(txn, username, ts, &records)
			return err
		},
	)
	if err != nil {
		return "", "", err
	}

	s.emitChange(ctx, tokenKind, username, changefeed.OpWrite, ts.NodeID, records)
	return token, tokenUuid, nil
}

func (s *storageImpl) VerifyToken(ctx context.Context, token string) (*models.User, error) {
	var user *models.User

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			user, err = transactions.VerifyToken(txn, token)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storageImpl) ListTokens(ctx context.Context, username string) ([]string, error) {
	var tokens []string

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			tokens, err = transactions.ListTokens(txn, username)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *storageImpl) DeleteToken(ctx context.Context, username string, tokenUUID string) error {
	ts := s.clock.Now()
	var records []changefeed.Record
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteToken(txn, username, tokenUUID, ts, &records)
		},
	)
	if err != nil {
		return err
	}

	s.emitChange(ctx, tokenKind, username, changefeed.OpDelete, ts.NodeID, records)
	return nil
}

func (s *storageImpl) emitChange(
	ctx context.Context,
	kind string,
	name string,
	op changefeed.Operation,
	origin string,
	records []changefeed.Record,
) {
	event := changefeed.ChangeEvent{
		Namespace: changefeed.NamespaceAuth,
		Kind:      kind,
		Name:      name,
		Op:        op,
		Origin:    origin,
		Records:   records,
	}
	if err := s.notifier.Notify(ctx, event); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to emit change event",
			"channel", "storage.auth",
			"kind", kind,
			"name", name,
			"error", err.Error(),
		)
	}
}

func (s *storageImpl) emitResync(ctx context.Context) {
	event := changefeed.ChangeEvent{
		Namespace: changefeed.NamespaceAuth,
		Resync:    true,
	}
	if err := s.notifier.Notify(ctx, event); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to emit resync event",
			"channel", "storage.auth",
			"error", err.Error(),
		)
	}
}
