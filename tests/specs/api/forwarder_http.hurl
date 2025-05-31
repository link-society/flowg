###############################################################################
# DESCRIPTION: Test HTTP Forwarder
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Create HTTP Forwarder
#  - Send a record to the forwarder and verify it succeeds
#  - Delete the forwarder

PUT http://localhost:5080/api/v1/forwarders/http-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder": {
    "config": {
      "type": "http",
      "url": "http://test-flowg-mockserver:1080/test/flowg/forwarder/http/success",
      "headers": {
        "Foo": "Bar"
      }
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/http-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "Hello, World!"
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

DELETE http://localhost:5080/api/v1/forwarders/http-success
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# -----------------------------------------------------------------------------
# TEST:
#  - Create HTTP Forwarder with failure response
#  - Send a record to the forwarder and verify it fails
#  - Delete the forwarder

PUT http://localhost:5080/api/v1/forwarders/http-fail
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder": {
    "config": {
      "type": "http",
      "url": "http://test-flowg-mockserver:1080/test/flowg/forwarder/http/fail",
      "headers": {
        "Foo": "Bar"
      }
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/http-fail
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "Hello, World!"
  }
}
HTTP 500

DELETE http://localhost:5080/api/v1/forwarders/http-fail
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
