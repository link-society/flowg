###############################################################################
# DESCRIPTION: Test Saml auth provider
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Check Create Saml auth provider
#  - Verify the Saml auth provider appears in the list
#  - Verify the Saml auth provider can be fetched individually
#  - Delete the Saml auth provider
#  - Verify the Saml auth provider is deleted

PUT http://localhost:5080/api/v1/auth-providers/test-saml
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "auth_provider": {
    "name": "test-saml",
    "display_name": "Test SAML Connection",
    "config": {
      "type": "saml",
      "idp_metadata_url": "https://example.com"
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
jsonpath "$.auth_providers[0].name" == "test-saml"

GET http://localhost:5080/api/v1/auth-providers/test-saml
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.auth_provider.name" == "test-saml"
jsonpath "$.auth_provider.display_name" == "Test SAML Connection"

DELETE http://localhost:5080/api/v1/auth-providers/test-saml
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/auth-providers/test-saml
Authorization: Bearer {{admin_token}}
HTTP 404
