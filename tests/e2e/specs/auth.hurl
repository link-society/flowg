GET http://localhost:5080/api/v1/whoami
HTTP 401

GET http://localhost:5080/api/v1/whoami
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/whoami
Authorization: Bearer {{guest_token}}
HTTP 200
