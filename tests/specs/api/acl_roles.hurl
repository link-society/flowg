###############################################################################
# DESCRIPTION: Test role management through API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify the default admin role exists

GET http://localhost:5080/api/v1/roles
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.roles" count == 1
jsonpath "$.roles[*].name" contains "admin"

# -----------------------------------------------------------------------------
# TEST:
#  - Verify guest user cannot access ACLs
#  - Create a new role with permission to read ACLs
#  - Assign it to the guest user
#  - Verify guest user can now access ACLs

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{guest_token}}
HTTP 403

PUT http://localhost:5080/api/v1/roles/test
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "scopes": [
    "read_acls"
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/roles
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.roles" count == 2
jsonpath "$.roles[*].name" contains "test"

PUT http://localhost:5080/api/v1/users/guest
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "password": "guest",
  "roles": ["test"]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{guest_token}}
HTTP 200

# -----------------------------------------------------------------------------
# TEST:
#  - Remove the test role from the guest user
#  - Verify guest user can no longer access ACLs

PUT http://localhost:5080/api/v1/users/guest
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "password": "guest",
  "roles": []
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{guest_token}}
HTTP 403

# -----------------------------------------------------------------------------
# TEST:
#  - Delete the test role
#  - Verify it is no longer listed

DELETE http://localhost:5080/api/v1/roles/test
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/roles
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.roles" count == 1
jsonpath "$.roles[*].name" not contains "test"
