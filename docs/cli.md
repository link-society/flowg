# Command Line Interface

```
Low-Code log management solution

Usage:
  flowg [command]

Available Commands:
  admin       Admin commands (please run while the server is down)
  help        Help about any command
  serve       Start FlowG standalone server

Flags:
  -h, --help   help for flowg

Use "flowg [command] --help" for more information about a command.
```

## 1. `flowg serve`

```
Start FlowG standalone server

Usage:
  flowg serve [flags]

Flags:
      --auth-dir string     Path to the auth database directory (default "./data/auth")
      --bind string         Address to bind the server to (default ":5080")
      --config-dir string   Path to the config directory (default "./data/config")
  -h, --help                help for serve
      --log-dir string      Path to the log database directory (default "./data/logs")
      --verbose             Enable verbose logging
```

## 2. `flowg admin`

```
Admin commands (please run while the server is down)

Usage:
  flowg admin [command]

Available Commands:
  role        Role related admin commands (please run while the server is down)
  token       Personal Access Token related admin commands (please run while the server is down)
  user        User related admin commands (please run while the server is down)

Flags:
  -h, --help   help for admin

Use "flowg admin [command] --help" for more information about a command.
```

### 2.1. `flowg admin role`

```
Role related admin commands (please run while the server is down)

Usage:
  flowg admin role [command]

Available Commands:
  create      Create a new role
  delete      Delete an existing role
  list        List existing roles

Flags:
  -h, --help   help for role

Use "flowg admin role [command] --help" for more information about a command.
```

#### 2.1.1. `flowg admin role create`

```
Create a new role

Usage:
  flowg admin role create [flags]

Flags:
      --auth-dir string   Path to the log database directory (default "./data/auth")
  -h, --help              help for create
      --name string       Name of the role
```

#### 2.1.2. `flowg admin role delete`

```
Delete an existing role

Usage:
  flowg admin role delete [flags]

Flags:
      --auth-dir string   Path to the log database directory (default "./data/auth")
  -h, --help              help for delete
      --name string       Name of the role
```

#### 2.1.3. `flowg admin role list`

```
List existing roles

Usage:
  flowg admin role list [flags]

Flags:
      --auth-dir string   Path to the log database directory (default "./data/auth")
  -h, --help              help for list
```

### 2.2. `flowg admin user`

```
User related admin commands (please run while the server is down)

Usage:
  flowg admin user [command]

Available Commands:
  create      Create a new user
  delete      Delete an existing user
  list        List existing users

Flags:
  -h, --help   help for user

Use "flowg admin user [command] --help" for more information about a command.
```

#### 2.2.1. `flowg admin user create`

```
Create a new user

Usage:
  flowg admin user create [flags]

Flags:
      --auth-dir string   Path to the log database directory (default "./data/auth")
  -h, --help              help for create
      --name string       Name of the user
      --password string   Password of the user
```

#### 2.2.2. `flowg admin user delete`

```
Delete an existing user

Usage:
  flowg admin user delete [flags]

Flags:
      --auth-dir string   Path to the log database directory (default "./data/auth")
  -h, --help              help for delete
      --name string       Name of the user
```

#### 2.2.3. `flowg admin user list`

```
List existing users

Usage:
  flowg admin user list [flags]

Flags:
      --auth-dir string   Path to the log database directory (default "./data/auth")
  -h, --help              help for list
```

### 2.3. `flowg admin token`

```
Personal Access Token related admin commands (please run while the server is down)

Usage:
  flowg admin token [command]

Available Commands:
  create      Create a new Personal Access Token

Flags:
  -h, --help   help for token

Use "flowg admin token [command] --help" for more information about a command.
```

#### 2.3.1. `flowg admin token create`

```
Create a new Personal Access Token

Usage:
  flowg admin token create [flags]

Flags:
      --auth-dir string   Path to the log database directory (default "./data/auth")
  -h, --help              help for create
      --user string       Name of the user
```
