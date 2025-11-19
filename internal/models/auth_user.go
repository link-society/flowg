package models

type User struct {
	Name  string   `json:"name" required:"true" minLength:"1"`
	Roles []string `json:"roles" required:"true" items.minLength:"1"`
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role == roleName {
			return true
		}
	}

	return false
}
