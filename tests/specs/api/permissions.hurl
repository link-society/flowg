###############################################################################
# DESCRIPTION: Test access control
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that unauthenticated requests on streams are denied
#  - Verify that a user with the scope "read_streams" can access streams
#  - Verify that a user without the scope "read_streams" is denied access

GET http://localhost:5080/api/v1/streams
HTTP 401

GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{guest_token}}
HTTP 403

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that unauthenticated requests on transformers are denied
#  - Verify that a user with the scope "read_transformers" can access transformers
#  - Verify that a user without the scope "read_transformers" is denied access

GET http://localhost:5080/api/v1/transformers
HTTP 401

GET http://localhost:5080/api/v1/transformers
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/transformers
Authorization: Bearer {{guest_token}}
HTTP 403

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that unauthenticated requests on pipelines are denied
#  - Verify that a user with the scope "read_pipelines" can access pipelines
#  - Verify that a user without the scope "read_pipelines" is denied access

GET http://localhost:5080/api/v1/pipelines
HTTP 401

GET http://localhost:5080/api/v1/pipelines
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/pipelines
Authorization: Bearer {{guest_token}}
HTTP 403

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that unauthenticated requests on forwarders are denied
#  - Verify that a user with the scope "read_forwarders" can access forwarders
#  - Verify that a user without the scope "read_forwarders" is denied access

GET http://localhost:5080/api/v1/forwarders
HTTP 401

GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer {{guest_token}}
HTTP 403

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that unauthenticated requests on users are denied
#  - Verify that a user with the scope "read_acls" can access users
#  - Verify that a user without the scope "read_acls" is denied access

GET http://localhost:5080/api/v1/users
HTTP 401

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/users
Authorization: Bearer {{guest_token}}
HTTP 403

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that unauthenticated requests on roles are denied
#  - Verify that a user with the scope "read_acls" can access roles
#  - Verify that a user without the scope "read_acls" is denied access

GET http://localhost:5080/api/v1/roles
HTTP 401

GET http://localhost:5080/api/v1/roles
Authorization: Bearer {{admin_token}}
HTTP 200

GET http://localhost:5080/api/v1/roles
Authorization: Bearer {{guest_token}}
HTTP 403
