GET http://localhost:5080/api/v1/streams
HTTP 401

GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/transformers
HTTP 401

GET http://localhost:5080/api/v1/transformers
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/transformers
Authorization: Bearer {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/pipelines
HTTP 401

GET http://localhost:5080/api/v1/pipelines
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/pipelines
Authorization: Bearer {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/alerts
HTTP 401

GET http://localhost:5080/api/v1/alerts
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/alerts
Authorization: Bearer {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/users
HTTP 401

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/roles
HTTP 401

GET http://localhost:5080/api/v1/roles
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/roles
Authorization: Bearer {{guest_token}}
HTTP 403
