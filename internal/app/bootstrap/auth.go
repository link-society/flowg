package bootstrap

import (
	"fmt"

	"link-society.com/flowg/internal/data/auth"
)

func DefaultRolesAndUsers(authDb *auth.Database) error {
	roleSys := auth.NewRoleSystem(authDb)
	userSys := auth.NewUserSystem(authDb)

	roles, err := roleSys.ListRoles()
	if err != nil {
		return err
	}

	if len(roles) == 0 {
		adminRole := auth.Role{
			Name: "admin",
			Scopes: []auth.Scope{
				auth.SCOPE_SEND_LOGS,
				auth.SCOPE_WRITE_ACLS,
				auth.SCOPE_WRITE_PIPELINES,
				auth.SCOPE_WRITE_TRANSFORMERS,
				auth.SCOPE_WRITE_STREAMS,
				auth.SCOPE_WRITE_ALERTS,
			},
		}

		err := roleSys.SaveRole(adminRole)
		if err != nil {
			return fmt.Errorf("failed to bootstrap admin role: %w", err)
		}
	}

	users, err := userSys.ListUsers()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		rootUser := auth.User{
			Name:  "root",
			Roles: []string{"admin"},
		}

		err := userSys.SaveUser(rootUser, "root")
		if err != nil {
			return fmt.Errorf("failed to bootstrap root user: %w", err)
		}
	}

	return nil
}
