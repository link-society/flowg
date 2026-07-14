###############################################################################
# DESCRIPTION: Test Azure Monitor Forwarder
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Create Azure Monitor Forwarder
#  - Send a record to the forwarder and verify it succeeds
#  - Delete the forwarder

PUT http://localhost:5080/api/v1/forwarders/monitor-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder":{
    "config":{
      "type":"azuremonitor",
      "endpoint":"https://floci-az-test-server:4577",
      "token": "token",
      "expires_on": "2006-01-02T15:04:05Z",
      "rule_id": "{{azuremonitor_dcr}}",
      "stream_name":"flowg",
      "allow_insecure": true
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/monitor-success
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

DELETE http://localhost:5080/api/v1/forwarders/monitor-success
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
