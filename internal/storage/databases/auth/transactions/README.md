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
- `user:<name>:password` — Argon2id hash of the user's password.
- `user:<name>:role:<role>` — one key per role assigned to the user.

### Roles

- `index:role:<name>` — existence marker used to enumerate roles.
- `role:<name>:<scope>` — one key per scope granted to the role.

### Personal access tokens

- `pat:<name>:<uuid>` — one key per token; its value is the SHA-256 hash of the
  token. A fast hash is enough because the token is a long random secret, not a
  low-entropy password.
- `index:pat:<sha256(token)>` — reverse index mapping a token's hash back to its
  owning `<name>`. The plaintext is returned to the caller only once, at
  creation; verification hashes the presented token and resolves the owner with a
  single O(1) lookup on this index (no per-token scan).

## Notes

- Writing a user or role reconciles its child keys: existing keys are diffed
  against the desired set, the missing ones are created and the obsolete ones are
  deleted.
- Authorization resolves a user's effective scopes by walking
  `user:<name>:role:*` and then `role:<role>:*`, with each write scope implying
  its matching read scope.
