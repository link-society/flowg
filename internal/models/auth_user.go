package models

type User struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role == roleName {
			return true
		}
	}

	return false
}
