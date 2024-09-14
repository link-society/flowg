GET http://localhost:5080/api/v1/auth/whoami
HTTP 401

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Bearer {{guest_token}}
HTTP 200

POST http://localhost:5080/api/v1/auth/login
Content-Type: application/json
{
  "username": "root",
  "password": "root"
}
HTTP 200
[Captures]
session_id: cookie "session_id"
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/auth/whoami
[Cookies]
session_id: {{session_id}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
