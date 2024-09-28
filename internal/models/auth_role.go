package models

type Role struct {
	Name   string  `json:"name"`
	Scopes []Scope `json:"scopes"`
}

func (r Role) HasScope(scope Scope) bool {
	for _, s := range r.Scopes {
		if s == scope {
			return true
		}
	}

	return false
}
