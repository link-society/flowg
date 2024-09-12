GET http://localhost:5080/api/v1/whoami
HTTP 401

GET http://localhost:5080/api/v1/whoami
Authorization: Bearer {{token}}
HTTP 200
