###############################################################################
# DESCRIPTION: Test user management through API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify the default admin and guest users exist

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.users" count == 2
jsonpath "$.users[*].name" contains "root"
jsonpath "$.users[*].name" contains "guest"
jsonpath "$.users[?(@.name == 'root')].roles[*]" count == 1
jsonpath "$..users[?(@.name == 'root')].roles[*]" contains "admin"
jsonpath "$.users[?(@.name == 'guest')].roles[*]" count == 0

# -----------------------------------------------------------------------------
# TEST:
#  - Create a new user with no roles
#  - Verify it appears in the list of users

PUT http://localhost:5080/api/v1/users/test
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "password": "test",
  "roles": []
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.users" count == 3
jsonpath "$.users[*].name" contains "test"

# -----------------------------------------------------------------------------
# TEST:
#  - Delete the new user
#  - Verify it no longer appears in the list of users

DELETE http://localhost:5080/api/v1/users/test
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.users" count == 2
jsonpath "$.users[*].name" not contains "test"
