###############################################################################
# DESCRIPTION: Test NATS JetStream Forwarder
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Create NATS JetStream Forwarder
#  - Send a record to the forwarder and verify it succeeds
#  - Delete the forwarder

PUT http://localhost:5080/api/v1/forwarders/jetstream-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder": {
    "config": {
      "type": "jetstream",
      "auth_mode": "none",
      "body": "@expr:toJSON(log)",
      "servers": [
          "http://nats-test-server:4222"
      ],
      "subject": "logs",
      "tls_insecure_skip_verify": true
    }
  }
}

HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/jetstream-success
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

DELETE http://localhost:5080/api/v1/forwarders/jetstream-success
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
