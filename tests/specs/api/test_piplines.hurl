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
  "flow": {
    "nodes": [
      {
        "id": "__builtin__source_direct",
        "type": "source",
        "position": {"x": 210, "y": 195},
        "deletable": false,
        "data": {"type": "direct"},
        "measured": {"width": 136, "height": 38},
        "selected": true,
        "dragging": false
      },
      {
        "id": "__builtin__source_syslog",
        "type": "source",
        "position": {"x": 210, "y": 250},
        "deletable": false,
        "data": {"type": "syslog"},
        "measured": {"width": 136, "height": 38},
        "selected": true,
        "dragging": false
      },
      {
        "id": "node-1",
        "type": "router",
        "position": {"x": 405, "y": 195},
        "data": {"stream": "test"},
        "measured": {"width": 241,"height": 91},
        "selected": false,
        "dragging": false
      }
    ],
    "edges": [
      {
        "id": "xy-edge____builtin__source_direct-node-1",
        "type": "smoothstep",
        "source": "__builtin__source_direct",
        "sourceHandle": "",
        "target": "node-1",
        "animated": true
      },
      {
        "id": "xy-edge____builtin__source_syslog-node-1",
        "type": "smoothstep",
        "source": "__builtin__source_syslog",
        "sourceHandle": "",
        "target": "node-1",
        "animated": true
      }
    ]
  },
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
