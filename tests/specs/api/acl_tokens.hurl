###############################################################################
# DESCRIPTION: Test token management through API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that the default admin user has one token

GET http://localhost:5080/api/v1/tokens
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.token_uuids" count == 1

# -----------------------------------------------------------------------------
# TEST:
#  - Create a new admin token
#  - Verify it appears in the list of tokens
#  - Verify it can be used to authenticate

POST http://localhost:5080/api/v1/token
Authorization: Bearer {{admin_token}}
HTTP 200
[Captures]
new_admin_token: jsonpath "$.token"
new_admin_token_uuid: jsonpath "$.token_uuid"
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/tokens
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.token_uuids" count == 2
jsonpath "$.token_uuids" contains "{{new_admin_token_uuid}}"

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer {{new_admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.user.name" == "root"

# -----------------------------------------------------------------------------
# TEST:
#  - Delete the new admin token
#  - Verify it can no longer be used to authenticate

DELETE http://localhost:5080/api/v1/tokens/{{new_admin_token_uuid}}
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer {{new_admin_token}}
HTTP 401
