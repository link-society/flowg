###############################################################################
# OTLP Forwarder - Success Scenario
###############################################################################

# Create OTLP Success Forwarder
PUT http://localhost:5080/api/v1/forwarders/otlp-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "forwarder": {
    "version": 2,
    "config": {
      "type": "otlp",
      "config": {
        "endpoint": "http://test-flowg-mockserver:1080/v1/logs",
        "headers": {
          "X-Custom-Header": "test-value",
          "X-Test-Mode": "success",
          "Content-Type": "application/x-protobuf"
        }
      }
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Send minimal test log (should forward with above headers)
POST http://localhost:5080/api/v1/test/forwarders/otlp-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "record": {
    "timestamp": "1700000000000000000",
    "body": "Minimal test log message"
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Send full test log
POST http://localhost:5080/api/v1/test/forwarders/otlp-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "record": {
    "timestamp": "1700000000000000000",
    "body": "Full test log message",
    "severity_number": "9",
    "severity_text": "INFO",
    "service": "auth",
    "environment": "production"
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Delete OTLP Success Forwarder
DELETE http://localhost:5080/api/v1/forwarders/otlp-success
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true


###############################################################################
# OTLP Forwarder - Failure Scenario
###############################################################################

# Create OTLP Fail Forwarder (headers specify X-Test-Mode: fail)
PUT http://localhost:5080/api/v1/forwarders/otlp-fail
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "forwarder": {
    "version": 2,
    "config": {
      "type": "otlp",
      "config": {
        "endpoint": "http://test-flowg-mockserver:1080/v1/logs",
        "headers": {
          "X-Custom-Header": "test-value",
          "X-Test-Mode": "fail",
          "Content-Type": "application/x-protobuf"
        }
      }
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Send test log that should fail (mockserver returns 500)
POST http://localhost:5080/api/v1/test/forwarders/otlp-fail
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "record": {
    "timestamp": "1700000000000000000",
    "body": "Test OTLP failure"
  }
}
HTTP 500
[Asserts]
jsonpath "$.status" == "INTERNAL"
jsonpath "$.error" == "internal: unexpected status code: 500, body: Internal Server Error"

# Delete OTLP Fail Forwarder
DELETE http://localhost:5080/api/v1/forwarders/otlp-fail
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
