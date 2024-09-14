GET http://localhost:5080/api/v1/tokens
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.token-uuids" count == 1

POST http://localhost:5080/api/v1/token
Authorization: Basic {{admin_token}}
HTTP 200
[Captures]
new_admin_token: jsonpath "$.token"
new_admin_token_uuid: jsonpath "$.token-uuid"
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/tokens
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.token-uuids" count == 2
jsonpath "$.token-uuids" includes "{{new_admin_token_uuid}}"

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Basic {{new_admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.user.name" == "root"

DELETE http://localhost:5080/api/v1/tokens/{{new_admin_token_uuid}}
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/auth/whoami
Authorization: Basic {{new_admin_token}}
HTTP 401
