###############################################################################
# DESCRIPTION: Test forwarder management through API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that no forwarders exist initially
#  - Verify that trying to get a non-existent forwarder returns 404

GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarders" count == 0

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
HTTP 404

# -----------------------------------------------------------------------------
# TEST:
#  - Create a forwarder
#  - Verify the forwarder appears in the list
#  - Verify the forwarder can be fetched individually

PUT http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder": {
    "config": {
      "type": "http",
      "url": "http://httpbin.org/anything",
      "headers": {
        "Foo": "Bar"
      }
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarders" count == 1
jsonpath "$.forwarders[0]" == "httpbin"

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarder.config.type" == "http"
jsonpath "$.forwarder.config.url" == "http://httpbin.org/anything"
jsonpath "$.forwarder.config.headers.Foo" == "Bar"

# -----------------------------------------------------------------------------
# TEST:
#  - Delete the forwarder
#  - Verify the forwarder is no longer listed

DELETE http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
HTTP 404
