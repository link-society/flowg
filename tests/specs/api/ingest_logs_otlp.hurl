###############################################################################
# DESCRIPTION: Test OpenTelemetry Logs Ingestion
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Send OpenTelemetry logs in JSON format to the default pipeline
#  - Send OpenTelemetry logs in Protobuf format to the default pipeline
#  - Verify logs have been ingested
#  - Purge the stream

POST http://localhost:5080/api/v1/pipelines/default/logs/otlp
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "resourceLogs": [{
    "resource": {
      "attributes": [{
        "key": "service.name",
        "value": { "stringValue": "demo-service" }
      }]
    },
    "scopeLogs": [{
      "scope": {
        "name": "demo-logger"
      },
      "logRecords": [{
        "timeUnixNano": 1713300000000000000,
        "severityNumber": 9,
        "severityText": "INFO",
        "body": { "stringValue": "Hello from Hurl" },
        "attributes": [{
          "key": "env",
          "value": { "stringValue": "dev" }
        }]
      }]
    }]
  }]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

POST http://localhost:5080/api/v1/pipelines/default/logs/otlp
Authorization: Bearer {{admin_token}}
Content-Type: application/x-protobuf
file,ingest_otlp.pb;
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 2
jsonpath "$.records[0].fields.body" == "Hello from Python"
jsonpath "$.records[1].fields.body" == "Hello from Hurl"

DELETE http://localhost:5080/api/v1/streams/default
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
