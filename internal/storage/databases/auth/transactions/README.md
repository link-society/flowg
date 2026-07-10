# transactions

The package at `internal/storage/databases/auth/transactions` holds the
low-level read/write operations behind the auth storage. Each function takes a
transaction and manipulates the auth key space directly; this document describes
how that key space is laid out.

## Key space

Users, roles and tokens are not stored as single serialized documents. Each fact
is its own key, so a set — a user's roles, a role's scopes — is simply the
collection of keys that exist under a given prefix.

### Users

- `index:user:<name>` — existence marker used to enumerate users.
- `user:<name>:password` — bcrypt hash of the user's password.
- `user:<name>:role:<role>` — one key per role assigned to the user.

### Roles

- `index:role:<name>` — existence marker used to enumerate roles.
- `role:<name>:<scope>` — one key per scope granted to the role.

### Personal access tokens

- `pat:<name>:<uuid>` — one key per token; its value is the bcrypt hash of the
  token. The plaintext is returned to the caller only once, at creation, so
  verification re-hashes the presented token and compares it against the stored
  hashes.

## Notes

- Writing a user or role reconciles its child keys: existing keys are diffed
  against the desired set, the missing ones are created and the obsolete ones are
  deleted.
- Authorization resolves a user's effective scopes by walking
  `user:<name>:role:*` and then `role:<role>:*`, with each write scope implying
  its matching read scope.
