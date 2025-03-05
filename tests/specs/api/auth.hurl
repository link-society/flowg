GET http://localhost:5080/api/v1/auth/whoami
HTTP 401

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer pat:{{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.user.name" == "root"

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer pat:{{guest_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.user.name" == "guest"

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
Authorization: Bearer jwt:{{jwt}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.user.name" == "root"

POST http://localhost:5080/api/v1/auth/change-password
Authorization: Bearer pat:{{guest_token}}
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
