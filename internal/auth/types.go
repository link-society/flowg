package auth

type Scope string

type Role struct {
	Name   string
	Scopes []Scope
}

type User struct {
	Name  string
	Roles []Role
}
