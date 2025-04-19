POST http://localhost:5080/api/v1/pipelines/default/logs/otlp
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "resource_logs": [{
    "resource": {
      "attributes": [{
        "key": "service.name",
        "value": { "stringValue": "demo-service" }
      }]
    },
    "scope_logs": [{
      "scope": {
        "name": "demo-logger"
      },
      "log_records": [{
        "time_unix_nano": 1713300000000000000,
        "severity_number": 9,
        "severity_text": "INFO",
        "body": { "stringValue": "Hello depuis curl !" },
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
