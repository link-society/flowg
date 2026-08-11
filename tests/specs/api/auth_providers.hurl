###############################################################################
# DESCRIPTION: Test auth provider management through API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that no auth providers exist initially
#  - Verify that trying to get a non-existent auth provider returns 404

GET http://localhost:5080/api/v1/auth-providers/
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.auth_providers" count == 0

GET http://localhost:5080/api/v1/auth-providers/test-oidc
Authorization: Bearer {{admin_token}}
HTTP 404

# -----------------------------------------------------------------------------
# TEST:
#  - Check Create auth provider endpoint is authenticated
#  - Create a auth provider
#  - Verify the auth provider appears in the list
#  - Check Read auth provider endpoint is authenticated
#  - Verify the auth provider can be fetched individually

PUT http://localhost:5080/api/v1/auth-providers/test-oidc
Content-Type: application/json
{
  "auth_provider": {
    "name": "test-oidc",
    "display_name": "Test OIDC Connection",
    "config": {
      "type": "oidc",
      "issuer": "https://example.com",
      "client_id": "abc",
      "client_secret": "123"
    }
  }
}
HTTP 401

PUT http://localhost:5080/api/v1/auth-providers/test-oidc
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "auth_provider": {
    "name": "test-oidc",
    "display_name": "Test OIDC Connection",
    "config": {
      "type": "oidc",
      "issuer": "https://example.com",
      "client_id": "abc",
      "client_secret": "123"
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/auth-providers/
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.auth_providers" count == 1
jsonpath "$.auth_providers[0].name" == "test-oidc"

GET http://localhost:5080/api/v1/auth-providers/test-oidc
HTTP 401

GET http://localhost:5080/api/v1/auth-providers/test-oidc
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.auth_provider.name" == "test-oidc"
jsonpath "$.auth_provider.display_name" == "Test OIDC Connection"

# -----------------------------------------------------------------------------
# TEST:
#  - Check Delete auth provider endpoint is authenticated
#  - Delete the auth provider
#  - Verify the auth provider is deleted

DELETE http://localhost:5080/api/v1/auth-providers/test-oidc
HTTP 401

DELETE http://localhost:5080/api/v1/auth-providers/test-oidc
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/auth-providers/test-oidc
Authorization: Bearer {{admin_token}}
HTTP 404
