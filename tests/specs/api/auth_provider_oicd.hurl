###############################################################################
# DESCRIPTION: Test Oicd auth provider
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Check Create Oicd auth provider
#  - Verify the Oicd auth provider appears in the list
#  - Verify the Oicd auth provider can be fetched individually
#  - Delete the Oicd auth provider
#  - Verify the Oicd auth provider is deleted

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
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.auth_provider.name" == "test-oidc"
jsonpath "$.auth_provider.display_name" == "Test OIDC Connection"

DELETE http://localhost:5080/api/v1/auth-providers/test-oidc
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/auth-providers/test-oidc
Authorization: Bearer {{admin_token}}
HTTP 404
