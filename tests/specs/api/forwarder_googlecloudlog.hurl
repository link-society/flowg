###############################################################################
# DESCRIPTION: Test Google Cloud Log Forwarder
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Create Google Cloud Log Forwarder
#  - Send a record to the forwarder and verify it succeeds
#  - Delete the forwarder

PUT http://localhost:5080/api/v1/forwarders/google-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder":{
    "config":{
      "type":"googlecloudlogging",
      "project_id":"flowg",
      "log_id":"flowg",
      "host":"floci-gcp-test-server",
      "port":4588,
      "disable_tls": true,
      "disable_auth": true
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/google-success
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

DELETE http://localhost:5080/api/v1/forwarders/google-success
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
