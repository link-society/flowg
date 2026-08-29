###############################################################################
# DESCRIPTION: Test that pipeline lifecycle changes propagate across a 3-node
#              FlowG cluster through the FoundationDB changelog watcher: a
#              pipeline created (or deleted) on one node must eventually be
#              running (or terminated) on every node, including the writer,
#              which observes its own change through the watch.
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Create a pipeline "watched" on node0, routing to the stream "watched"
#  - Verify every node eventually runs it (ingestion succeeds)
#  - Verify the ingested records all landed in the stream

PUT {{node0_url}}/api/v1/pipelines/watched
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "flow": {
    "nodes": [
      {
        "id": "__builtin__source_direct",
        "type": "source",
        "position": {"x": 210, "y": 195},
        "data": {"type": "direct"}
      },
      {
        "id": "node-router",
        "type": "router",
        "position": {"x": 405, "y": 195},
        "data": {"stream": "watched"}
      }
    ],
    "edges": [
      {
        "id": "xy-edge____builtin__source_direct-node-router",
        "source": "__builtin__source_direct",
        "sourceHandle": "",
        "target": "node-router"
      }
    ]
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# the changelog watch fires asynchronously: retry until each node runs the
# pipeline (ingestion returns 404 until then)
POST {{node0_url}}/api/v1/pipelines/watched/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
[Options]
retry: 30
retry-interval: 1000ms
{
  "records": [
    { "message": "watched-from-node0" }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

POST {{node1_url}}/api/v1/pipelines/watched/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
[Options]
retry: 30
retry-interval: 1000ms
{
  "records": [
    { "message": "watched-from-node1" }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

POST {{node2_url}}/api/v1/pipelines/watched/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
[Options]
retry: 30
retry-interval: 1000ms
{
  "records": [
    { "message": "watched-from-node2" }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

# the records went through the new pipeline into its stream
GET {{node1_url}}/api/v1/streams/watched/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 3
jsonpath "$.records[*].fields.message" contains "watched-from-node0"
jsonpath "$.records[*].fields.message" contains "watched-from-node1"
jsonpath "$.records[*].fields.message" contains "watched-from-node2"

# -----------------------------------------------------------------------------
# TEST:
#  - Delete the pipeline from node1
#  - Verify every node eventually terminates it (ingestion returns 404)

DELETE {{node1_url}}/api/v1/pipelines/watched
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST {{node0_url}}/api/v1/pipelines/watched/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
[Options]
retry: 30
retry-interval: 1000ms
{
  "records": [
    { "message": "after-delete" }
  ]
}
HTTP 404

POST {{node2_url}}/api/v1/pipelines/watched/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
[Options]
retry: 30
retry-interval: 1000ms
{
  "records": [
    { "message": "after-delete" }
  ]
}
HTTP 404
