GET http://localhost:5080/api/v1/users
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.users" count == 2
jsonpath "$.users[*].name" includes "root"
jsonpath "$.users[*].name" includes "guest"
jsonpath "$.users[?(@.name == 'root')].roles[*]" count == 1
jsonpath "$..users[?(@.name == 'root')].roles[*]" includes "admin"
jsonpath "$.users[?(@.name == 'guest')].roles[*]" count == 0

PUT http://localhost:5080/api/v1/users/test
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "password": "test",
  "roles": []
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/users
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.users" count == 3
jsonpath "$.users[*].name" includes "test"

DELETE http://localhost:5080/api/v1/users/test
Authorization: Basic {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/users
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.users" count == 2
jsonpath "$.users[*].name" not includes "test"
