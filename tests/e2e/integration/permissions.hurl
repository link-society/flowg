GET http://localhost:5080/api/v1/streams
HTTP 401

GET http://localhost:5080/api/v1/streams
Authorization: Basic {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/streams
Authorization: Basic {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/transformers
HTTP 401

GET http://localhost:5080/api/v1/transformers
Authorization: Basic {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/transformers
Authorization: Basic {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/pipelines
HTTP 401

GET http://localhost:5080/api/v1/pipelines
Authorization: Basic {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/pipelines
Authorization: Basic {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/alerts
HTTP 401

GET http://localhost:5080/api/v1/alerts
Authorization: Basic {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/alerts
Authorization: Basic {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/users
HTTP 401

GET http://localhost:5080/api/v1/users
Authorization: Basic {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/users
Authorization: Basic {{guest_token}}
HTTP 403

GET http://localhost:5080/api/v1/roles
HTTP 401

GET http://localhost:5080/api/v1/roles
Authorization: Basic {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/roles
Authorization: Basic {{guest_token}}
HTTP 403
