---
sidebar_position: 6
---

# Database Backup

**FlowG** provides CLI commands to perform **offline** (aka: while FlowG is
**not** running) backup and restore.

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
 - `./backup/config/`: containing a copy of the configuration folder

No write is ever done on the original databases.

## Performing an offline restore

```bash
flowg admin restore \
  --auth-dir ./data/auth \
  --config-dir ./data/config \
  --log-dir ./data/logs \
  --backup-dir ./backup
```

This command will expect the following files and directories to exist:

 - `./backup/auth.db`: containing the full snapshot of the authentication database
 - `./backup/log.db`: containing the full snapshot of the logs database
 - `./backup/config/`: containing a copy of the configuration folder

The content of the `auth.db` file will be inserted in the authentication database,
overwriting pre-existing keys, but not deleting new ones.

The content of the `log.db` file will be inserted in the logs database,
overwriting pre-existing keys, but not deleting new ones.

The content of the `config/` folder will be copied in the configuration folder,
pre-existing files will be overwritten, but new ones won't be deleted.

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
