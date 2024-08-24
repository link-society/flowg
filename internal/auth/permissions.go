package auth

import "fmt"

const (
	SCOPE_READ_PIPELINES     Scope = "read_pipelines"
	SCOPE_WRITE_PIPELINES    Scope = "write_pipelines"
	SCOPE_READ_TRANSFORMERS  Scope = "read_transformers"
	SCOPE_WRITE_TRANSFORMERS Scope = "write_transformers"
	SCOPE_READ_STREAMS       Scope = "read_streams"
	SCOPE_WRITE_STREAMS      Scope = "write_streams"
	SCOPE_CREATE_USERS       Scope = "create_users"
	SCOPE_SEND_LOGS          Scope = "send_logs"
)

func ParseScope(s string) (Scope, error) {
	switch s {
	case "read_pipelines":
		return SCOPE_READ_PIPELINES, nil
	case "write_pipelines":
		return SCOPE_WRITE_PIPELINES, nil
	case "read_transformers":
		return SCOPE_READ_TRANSFORMERS, nil
	case "write_transformers":
		return SCOPE_WRITE_TRANSFORMERS, nil
	case "read_streams":
		return SCOPE_READ_STREAMS, nil
	case "write_streams":
		return SCOPE_WRITE_STREAMS, nil
	case "create_users":
		return SCOPE_CREATE_USERS, nil
	case "send_logs":
		return SCOPE_SEND_LOGS, nil
	default:
		return "", fmt.Errorf("invalid scope: %s", s)
	}
}

func (r Role) HasScope(scope Scope) bool {
	for _, s := range r.Scopes {
		if s == scope {
			return true
		}
	}

	return false
}

func (u *User) Can(scope Scope) bool {
	for _, role := range u.Roles {
		for _, s := range role.Scopes {
			if s == scope {
				return true
			}
		}
	}

	return false
}
