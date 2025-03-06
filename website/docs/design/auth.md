---
sidebar_position: 2
---

# How Authentication Works?

FlowG uses a *Role Based Access Control* design (RBAC for short). Each role
assigns permissions to one or more of the following **scopes**:

| Scope Name | Description |
| --- | --- |
| `read_pipelines` | Can see a pipeline's flow, but cannot update it nor delete it |
| `write_pipelines` | Can create, read, update or delete a pipeline flow |
| `read_transformers` | Can see the source code of a transformer, but cannot update it nor delete it |
| `write_transformers` | Can create, read, update, or delete a transformer script |
| `read_streams` | Can query a stream |
| `write_streams` | Can purge a stream |
| `read_alerts` | Can see alert webhooks |
| `write_alerts` | Can create, read, update or delete alert webhooks |
| `read_acls` | Can list users and roles, but cannot update them nor delete them |
| `write_acls`| Can create, read, update or delete roles and users |
| `send_logs` | Can send logs to a pipeline for processing (useful for log sources) |

Each user is associated to one or more roles. A user has a required password,
and can have zero or more personal access tokens.

## Password and Token encryption

Any secret is hashed using the [Argon2](https://en.wikipedia.org/wiki/Argon2)
algorithm.

## Storage

For each role, there will be an index key with the following format:

```
index:role:<role name>
```

For example:

```
index:role:admin
index:role:viewer
```

For each scope associated to the role, there will be a key with the following
format:

```
role:<role name>:<scope name>
```

For example:

```
role:admin:write_streams
role:admin:write_transformers
role:admin:write_pipelines
role:admin:write_acls
```

For each user, there will be an index key with the following format:

```
index:user:<username>
```

For example:

```
index:user:guest
```

Each user will have a key containing the hashed password with the following
format:

```
user:<username>:password = argon2(password)
```

For example:

```
user:admin:password = ...
```

For each role associated to the user, there will be a key with the following
format:

```
user:<username>:role:<role name>
```

For example:

```
user:guest:role:viewer
```

For each Personal Access Token associated to the user, there will be a key with
the following format:

```
pat:<username>:<uuid> = argon2(token)
```

For example:

```
pat:guest:f6c2424a-bc1f-4030-9e42-7a09b96452a7
```
