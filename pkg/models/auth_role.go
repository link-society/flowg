package models

// Role is a named bundle of Scopes. Users are granted roles, and a user's
// effective permissions are the union of the scopes of all its roles.
type Role struct {
	Name   string  `json:"name" required:"true" minLength:"1"`
	Scopes []Scope `json:"scopes" required:"true"`
}

// HasScope reports whether the role grants the given scope.
func (r Role) HasScope(scope Scope) bool {
	for _, s := range r.Scopes {
		if s == scope {
			return true
		}
	}

	return false
}
