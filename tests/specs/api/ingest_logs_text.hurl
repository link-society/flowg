###############################################################################
# DESCRIPTION: Test Textual Logs Ingestion
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Send textual logs to the default pipeline
#  - Verify logs have been ingested
#  - Purge the stream

POST http://localhost:5080/api/v1/pipelines/default/logs/text
Authorization: Bearer {{admin_token}}
Content-Type: text/plain
```
hello
world
```
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
jsonpath "$.records[0].fields.content" == "world"
jsonpath "$.records[1].fields.content" == "hello"

DELETE http://localhost:5080/api/v1/streams/default
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
