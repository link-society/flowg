###############################################################################
# DESCRIPTION: Test CloudWatch Forwarder
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Create CloudWatch Forwarder
#  - Send a record to the forwarder and verify it succeeds
#  - Delete the forwarder

PUT http://localhost:5080/api/v1/forwarders/cloudwatch-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder":{
    "config":{
      "type":"cloudwatch",
      "app_id":"flowg",
      "endpoint":"http://floci-test-server:4566",
      "region":"us-east-1",
      "access_key_id":"test",
      "secret_access_key":"test",
      "session_token":"",
      "group":"flowg",
      "stream":"logs"
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/cloudwatch-success
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

DELETE http://localhost:5080/api/v1/forwarders/cloudwatch-success
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
