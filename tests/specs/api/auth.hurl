###############################################################################
# DESCRIPTION: Test authentication
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that we can't access the API without authentication
#  - Verify that a personal access token can be used to authenticate

GET http://localhost:5080/api/v1/auth/whoami
HTTP 401

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.user.name" == "root"

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer {{guest_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.user.name" == "guest"

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that we can get a temporary JSON Web Token using credentials
#  - Verify that the token can be used to authenticate

POST http://localhost:5080/api/v1/auth/login
Content-Type: application/json
{
  "username": "root",
  "password": "root"
}
HTTP 200
[Captures]
jwt: jsonpath "$.token"
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer {{jwt}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.user.name" == "root"

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that we can change the password of the guest user
#  - Verify that the guest user can no longer authenticate with the old password
#  - Verify that the guest user can authenticate with the new password

POST http://localhost:5080/api/v1/auth/change-password
Authorization: Bearer {{guest_token}}
Content-Type: application/json
{
  "old_password": "guest",
  "new_password": "guest123"
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/auth/login
Content-Type: application/json
{
  "username": "guest",
  "password": "guest"
}
HTTP 401

POST http://localhost:5080/api/v1/auth/login
Content-Type: application/json
{
  "username": "guest",
  "password": "guest123"
}
HTTP 200
