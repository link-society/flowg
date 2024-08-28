package auth

type Scope string

type Role struct {
	Name   string  `json:"name"`
	Scopes []Scope `json:"scopes"`
}

type User struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func (r Role) HasScope(scope Scope) bool {
	for _, s := range r.Scopes {
		if s == scope {
			return true
		}
	}

	return false
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role == roleName {
			return true
		}
	}

	return false
}
