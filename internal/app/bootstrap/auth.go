package bootstrap

import (
	"context"
	"fmt"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/auth"
)

func DefaultRolesAndUsers(ctx context.Context, authStorage *auth.Storage) error {
	roles, err := authStorage.ListRoles(ctx)
	if err != nil {
		return err
	}

	if len(roles) == 0 {
		adminRole := models.Role{
			Name: "admin",
			Scopes: []models.Scope{
				models.SCOPE_SEND_LOGS,
				models.SCOPE_WRITE_ACLS,
				models.SCOPE_WRITE_PIPELINES,
				models.SCOPE_WRITE_TRANSFORMERS,
				models.SCOPE_WRITE_STREAMS,
				models.SCOPE_WRITE_FORWARDERS,
			},
		}

		err := authStorage.SaveRole(ctx, adminRole)
		if err != nil {
			return fmt.Errorf("failed to bootstrap admin role: %w", err)
		}
	}

	users, err := authStorage.ListUsers(ctx)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		rootUser := models.User{
			Name:  "root",
			Roles: []string{"admin"},
		}

		err := authStorage.SaveUser(ctx, rootUser, "root")
		if err != nil {
			return fmt.Errorf("failed to bootstrap root user: %w", err)
		}
	}

	return nil
}
