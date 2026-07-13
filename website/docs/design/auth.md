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
| `read_forwarders` | Can see forwarder configurations |
| `write_forwarders` | Can create, read, update or delete forwarder configurations |
| `read_acls` | Can list users and roles, but cannot update them nor delete them |
| `write_acls`| Can create, read, update or delete roles and users |
| `send_logs` | Can send logs to a pipeline for processing (useful for log sources) |
| `read_system_configuration` | Can see the global system configuration |
| `write_system_configuration` | Can edit the global system configuration |

Each user is associated to one or more roles. A user has a required password,
and can have zero or more personal access tokens.

## Password and Token hashing

User **passwords** are low-entropy, human-chosen secrets, so they are hashed with
the memory-hard [Argon2id](https://en.wikipedia.org/wiki/Argon2) key-derivation
function to make offline brute-forcing expensive.

**Personal Access Tokens** are long, randomly-generated high-entropy secrets that
cannot be feasibly brute-forced regardless of hash speed, so they are hashed with
a fast `SHA-256`.

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
user:<username>:password = argon2id(password)
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
pat:<username>:<uuid> = sha256(token)
```

For example:

```
pat:guest:f6c2424a-bc1f-4030-9e42-7a09b96452a7
```

Additionally, a reverse index lets a token be resolved to its owner in constant
time, without scanning every token:

```
index:pat:<sha256(token)> = <username>
```
