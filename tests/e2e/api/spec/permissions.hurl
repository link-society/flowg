GET http://localhost:5080/api/v1/streams
HTTP 401

GET http://localhost:5080/api/v1/streams
Authorization: Bearer pat:{{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/streams
Authorization: Bearer pat:{{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/transformers
HTTP 401

GET http://localhost:5080/api/v1/transformers
Authorization: Bearer pat:{{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/transformers
Authorization: Bearer pat:{{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/pipelines
HTTP 401

GET http://localhost:5080/api/v1/pipelines
Authorization: Bearer pat:{{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/pipelines
Authorization: Bearer pat:{{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/alerts
HTTP 401

GET http://localhost:5080/api/v1/alerts
Authorization: Bearer pat:{{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/alerts
Authorization: Bearer pat:{{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/users
HTTP 401

GET http://localhost:5080/api/v1/users
Authorization: Bearer pat:{{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/users
Authorization: Bearer pat:{{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/roles
HTTP 401

GET http://localhost:5080/api/v1/roles
Authorization: Bearer pat:{{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/roles
Authorization: Bearer pat:{{guest_token}}
HTTP 403
