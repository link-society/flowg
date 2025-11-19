package models

type Role struct {
	Name   string  `json:"name" required:"true" minLength:"1"`
	Scopes []Scope `json:"scopes" required:"true"`
}

func (r Role) HasScope(scope Scope) bool {
	for _, s := range r.Scopes {
		if s == scope {
			return true
		}
	}

	return false
}
