###############################################################################
# DESCRIPTION: The FoundationDB backend does not support snapshot backup and
#              restore. Every backup and restore endpoint must therefore fail
#              cleanly with HTTP 501 Not Implemented on a clustered deployment,
#              never a false-success 200 with an empty file.
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Every backup endpoint returns 501 Not Implemented (checked across nodes)

GET {{node0_url}}/api/v1/backup/auth
Authorization: Bearer {{admin_token}}
HTTP 501

GET {{node1_url}}/api/v1/backup/config
Authorization: Bearer {{admin_token}}
HTTP 501

GET {{node2_url}}/api/v1/backup/logs
Authorization: Bearer {{admin_token}}
HTTP 501

# -----------------------------------------------------------------------------
# TEST:
#  - Every restore endpoint returns 501 Not Implemented (checked across nodes)

POST {{node0_url}}/api/v1/restore/auth
Authorization: Bearer {{admin_token}}
[Multipart]
backup: file,fixtures/dummy.db;
HTTP 501

POST {{node1_url}}/api/v1/restore/config
Authorization: Bearer {{admin_token}}
[Multipart]
backup: file,fixtures/dummy.db;
HTTP 501

POST {{node2_url}}/api/v1/restore/logs
Authorization: Bearer {{admin_token}}
[Multipart]
backup: file,fixtures/dummy.db;
HTTP 501
