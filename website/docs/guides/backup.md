---
sidebar_position: 6
---

# Database Backup

## Performing a backup

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

## Performing a restore


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
  -H "Authorization: Bearer jwt:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/auth \
  -X POST --form backup=auth.db

curl \
  -H "Authorization: Bearer jwt:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/config \
  -X POST --form backup=config.db

curl \
  -H "Authorization: Bearer jwt:${FLOWG_TOKEN}" \
  http://localhost:5080/api/v1/restore/logs \
  -X POST --form backup=logs.db
```

Performing an online restore requires the following permissions:

 - auth: `Write ACLs`
 - config: `Write Pipelines`, `Write Transformers`, `Write Alerts`
 - logs: `Write Streams`
