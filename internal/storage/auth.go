package storage

import (
	"context"

	"link-society.com/flowg/internal/models"
)

// AuthStorage is the contract for persisting FlowG's authentication and
// authorization state: roles, users, their permission scopes, and the personal
// access tokens used to authenticate API calls.
//
// It embeds [Streamable] so the whole store can be backed up and restored as a
// stream. Implementations live under internal/storage/backends.
type AuthStorage interface {
	Streamable

	// ListRoles returns every role known to the store.
	ListRoles(ctx context.Context) ([]models.Role, error)
	// FetchRole returns the role with the given name, or an error if it does
	// not exist.
	FetchRole(ctx context.Context, name string) (*models.Role, error)
	// SaveRole creates or replaces a role.
	SaveRole(ctx context.Context, role models.Role) error
	// DeleteRole removes the role with the given name.
	DeleteRole(ctx context.Context, name string) error

	// ListUsers returns every user account known to the store.
	ListUsers(ctx context.Context) ([]models.User, error)
	// FetchUser returns the user with the given name, or an error if it does
	// not exist.
	FetchUser(ctx context.Context, name string) (*models.User, error)
	// ListUserScopes returns the permission scopes granted to the named user
	// through its roles.
	ListUserScopes(ctx context.Context, name string) ([]models.Scope, error)
	// SaveUser creates or replaces a user, setting the given password.
	SaveUser(ctx context.Context, user models.User, password string) error
	// PatchUserRoles updates only the role assignments of an existing user.
	PatchUserRoles(ctx context.Context, user models.User) error
	// DeleteUser removes the user with the given name.
	DeleteUser(ctx context.Context, name string) error

	// VerifyUserPassword reports whether the given password matches the stored
	// credentials of the named user.
	VerifyUserPassword(ctx context.Context, name, password string) (bool, error)
	// VerifyUserPermission reports whether the named user has been granted the
	// given permission scope.
	VerifyUserPermission(ctx context.Context, username string, scope models.Scope) (bool, error)

	// CreateToken mints a new personal access token for the named user and
	// returns the clear-text token together with its UUID.
	CreateToken(ctx context.Context, username string) (string, string, error)
	// VerifyToken resolves a clear-text personal access token to its owning
	// user, or returns an error if the token is unknown.
	VerifyToken(ctx context.Context, token string) (*models.User, error)
	// ListTokens returns the UUIDs of the personal access tokens owned by the
	// named user.
	ListTokens(ctx context.Context, username string) ([]string, error)
	// DeleteToken revokes the personal access token identified by tokenUUID for
	// the named user.
	DeleteToken(ctx context.Context, username string, tokenUUID string) error
}
