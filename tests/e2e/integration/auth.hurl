GET http://localhost:5080/api/v1/auth/whoami
HTTP 401

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Basic {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Basic {{guest_token}}
HTTP 200

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
