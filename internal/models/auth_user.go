package models

// User is an account identified by name and granted a set of roles (by name).
// The roles resolve to the scopes that make up the user's permissions.
type User struct {
	Name  string   `json:"name" required:"true" minLength:"1"`
	Roles []string `json:"roles" required:"true" items.minLength:"1"`
}

// HasRole reports whether the user is assigned the named role.
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role == roleName {
			return true
		}
	}

	return false
}
