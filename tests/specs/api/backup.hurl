GET http://localhost:5080/api/v1/backup/auth
Authorization: Bearer {{admin_token}}
[Options]
output: backup/auth.db
HTTP 200

GET http://localhost:5080/api/v1/backup/config
Authorization: Bearer {{admin_token}}
[Options]
output: backup/config.db
HTTP 200

GET http://localhost:5080/api/v1/backup/logs
Authorization: Bearer {{admin_token}}
[Options]
output: backup/logs.db
HTTP 200

POST http://localhost:5080/api/v1/restore/auth
Authorization: Bearer {{admin_token}}
[Multipart]
backup: file,backup/auth.db;
HTTP 200

POST http://localhost:5080/api/v1/restore/config
Authorization: Bearer {{admin_token}}
[Multipart]
backup: file,backup/config.db;
HTTP 200

POST http://localhost:5080/api/v1/restore/logs
Authorization: Bearer {{admin_token}}
[Multipart]
backup: file,backup/logs.db;
HTTP 200
