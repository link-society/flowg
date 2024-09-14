GET http://localhost:5080/api/v1/roles
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.roles" count == 1
jsonpath "$.roles[*].name" includes "admin"

PUT http://localhost:5080/api/v1/roles/test
Authorization: Basic {{admin_token}}
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
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.roles" count == 2
jsonpath "$.roles[*].name" includes "test"

GET http://localhost:5080/api/v1/users
Authorization: Basic {{guest_token}}
HTTP 403

PUT http://localhost:5080/api/v1/users/guest
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "password": "guest",
  "roles": ["test"]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/users
Authorization: Basic {{guest_token}}
HTTP 200

PUT http://localhost:5080/api/v1/users/guest
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "password": "guest",
  "roles": []
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/users
Authorization: Basic {{guest_token}}
HTTP 403

DELETE http://localhost:5080/api/v1/roles/test
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/roles
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.roles" count == 1
jsonpath "$.roles[*].name" not includes "test"
