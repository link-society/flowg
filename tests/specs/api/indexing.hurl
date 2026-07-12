###############################################################################
# DESCRIPTION: Test filtering stream queries by indexed field values
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Configure the "default" stream to index the "level" field
#  - Ingest records with different "level" values through the default pipeline
#  - Verify that querying without indexing returns every record in the window
#  - Verify that an indexed-field filter returns ONLY the matching records
#  - Verify OR-within-a-field semantics
#  - Verify an unknown value matches nothing
#  - Purge the stream

PUT http://localhost:5080/api/v1/streams/default
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "config": {
    "indexed_fields": ["level"],
    "ttl": 0,
    "size": 0
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/pipelines/default/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "records": [
    { "level": "error", "message": "boom" },
    { "level": "info", "message": "hello" },
    { "level": "info", "message": "world" }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 3

# Without indexing, every record in the time window is returned
GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 3

# Filtering on an indexed value returns ONLY the matching record
GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
indexing: {"level":["error"]}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 1
jsonpath "$.records[0].fields.level" == "error"
jsonpath "$.records[0].fields.message" == "boom"

# A different indexed value returns its own records
GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
indexing: {"level":["info"]}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 2

# OR within a field: both values match, so all records are returned
GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
indexing: {"level":["error","info"]}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 3

# A value that no record carries matches nothing
GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
indexing: {"level":["debug"]}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 0

DELETE http://localhost:5080/api/v1/streams/default
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
