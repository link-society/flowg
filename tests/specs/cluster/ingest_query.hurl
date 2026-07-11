###############################################################################
# DESCRIPTION: Test log ingestion and querying across a 3-node FlowG cluster
#              backed by a shared FoundationDB.
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Ingest a distinct log record on each of the 3 nodes
#  - Verify that every node can query all 3 records (shared FoundationDB store)

POST {{node0_url}}/api/v1/pipelines/default/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "records": [
    { "message": "from-node0" }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

POST {{node1_url}}/api/v1/pipelines/default/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "records": [
    { "message": "from-node1" }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

POST {{node2_url}}/api/v1/pipelines/default/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "records": [
    { "message": "from-node2" }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

# Query node0: all 3 records must be visible
GET {{node0_url}}/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 3
jsonpath "$.records[*].fields.message" contains "from-node0"
jsonpath "$.records[*].fields.message" contains "from-node1"
jsonpath "$.records[*].fields.message" contains "from-node2"

# Query node1: all 3 records must be visible
GET {{node1_url}}/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 3
jsonpath "$.records[*].fields.message" contains "from-node0"
jsonpath "$.records[*].fields.message" contains "from-node1"
jsonpath "$.records[*].fields.message" contains "from-node2"

# Query node2: all 3 records must be visible
GET {{node2_url}}/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 3
jsonpath "$.records[*].fields.message" contains "from-node0"
jsonpath "$.records[*].fields.message" contains "from-node1"
jsonpath "$.records[*].fields.message" contains "from-node2"

# -----------------------------------------------------------------------------
# TEST:
#  - Purge the stream from one node
#  - Verify it is empty on all nodes

DELETE {{node0_url}}/api/v1/streams/default
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET {{node1_url}}/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 0
