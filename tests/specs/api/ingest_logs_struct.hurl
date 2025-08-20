###############################################################################
# DESCRIPTION: Test Structured Logs Ingestion
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Send structured logs to the default pipeline
#  - Verify logs have been ingested
#  - Purge the stream

POST http://localhost:5080/api/v1/pipelines/default/logs/struct
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
jsonpath "$.processed_count" == 2

GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 2
jsonpath "$.records[*].fields.message" contains "world"
jsonpath "$.records[*].fields.message" contains "hello"

DELETE http://localhost:5080/api/v1/streams/default
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
