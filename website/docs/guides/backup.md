---
sidebar_position: 6
---

# Database Backup

**FlowG** provides:

 - CLI commands to perform **offline** (aka: while FlowG is **not** running)
   backup and restore
 - API endpoints to perform **online** (aka: while FlowG is running) backup and
   restore

## Performing an offline backup

```bash
flowg admin backup \
  --auth-dir ./data/auth \
  --config-dir ./data/config \
  --log-dir ./data/logs \
  --backup-dir ./backup
```

This command will create the following files and directories:

 - `./backup/auth.db`: containing the full snapshot of the authentication database
 - `./backup/log.db`: containing the full snapshot of the logs database
 - `./backup/config.db`: containing the full snapshot of the config database

No write is ever done on the original databases.

## Performing an offline restore

```bash
flowg admin restore \
  --auth-dir ./data/auth \
  --config-dir ./data/config \
  --log-dir ./data/logs \
  --backup-dir ./backup
```

This command will expect the following files to exist:

 - `./backup/auth.db`: containing the full snapshot of the authentication database
 - `./backup/log.db`: containing the full snapshot of the logs database
 - `./backup/config.db`: containing the full snapshot of the config database

The content of the `auth.db` file will be inserted in the authentication database,
overwriting pre-existing keys, but not deleting new ones.

The content of the `log.db` file will be inserted in the logs database,
overwriting pre-existing keys, but not deleting new ones.

The content of the `config.db` file will be inserted in the config database,
overwriting pre-existing items, but not deleting new ones.

If you want a *destructive* restore (aka: remove new items that were not backed
up), you need to remove the old database first:

```bash
rm -rf ./data
flowg admin restore \
  --auth-dir ./data/auth \
  --config-dir ./data/config \
  --log-dir ./data/logs \
  --backup-dir ./backup
```

## Performing an online backup

Using a Personal Access Token:

```bash
export FLOWG_TOKEN="<your token>"

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/backup/auth \
  --output auth.db

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/backup/config \
  --output config.db

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/backup/logs \
  --output logs.db
```

Using a JSON Web Token:

```bash
export FLOWG_TOKEN=$(
  curl \
    http://localhost:5080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username": "<your username>", "password": "<your password>"}' \
    2>/dev/null \
  | jq -r .token
)

curl \
  -H "Authorization: Bearer jwt:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/backup/auth \
  --output auth.db

curl \
  -H "Authorization: Bearer jwt:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/backup/config \
  --output config.db

curl \
  -H "Authorization: Bearer jwt:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/backup/logs \
  --output logs.db
```

Performing an online backup requires the following permissions:

 - auth: `Read ACLs`
 - config: `Read Pipelines`, `Read Transformers`, `Read Alerts`
 - logs: `Read Streams`

## Performing an online restore


Using a Personal Access Token:

```bash
export FLOWG_TOKEN="<your token>"

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/auth \
  -X POST --form backup=auth.db

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/config \
  -X POST --form backup=config.db

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/logs \
  -X POST --form backup=logs.db
```

Using a JSON Web Token:

```bash
export FLOWG_TOKEN=$(
  curl \
    http://localhost:5080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username": "<your username>", "password": "<your password>"}' \
    2>/dev/null \
  | jq -r .token
)

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/auth \
  -X POST --form backup=auth.db

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/config \
  -X POST --form backup=config.db

curl \
  -H "Authorization: Bearer pat:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/logs \
  -X POST --form backup=logs.db
```

Performing an online restore requires the following permissions:

 - auth: `Write ACLs`
 - config: `Write Pipelines`, `Write Transformers`, `Write Alerts`
 - logs: `Write Streams`
