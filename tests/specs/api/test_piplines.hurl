###############################################################################
# DESCRIPTION: Test the Pipeline test endpoint
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Test sending structured logs to the default pipeline

POST http://localhost:5080/api/v1/test/pipeline/default
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "records": [
    {
      "level": "info",
      "message": "hello"
    },
    {
      "level": "info",
      "message": "world"
    }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.trace" count == 4
jsonpath "$.trace[*].input.message" contains "world"
jsonpath "$.trace[*].input.message" contains "hello"
