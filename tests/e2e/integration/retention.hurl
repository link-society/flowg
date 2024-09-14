PUT http://localhost:5080/api/v1/streams/default
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "config": {
    "indexed_fields": [],
    "ttl": 2,
    "size": 0
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true


POST http://localhost:5080/api/v1/pipelines/default/logs
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "record": {
    "level": "info",
    "message": "Hello, World!"
  }
}
HTTP 200

GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Basic {{admin_token}}
[QueryStringParams]
from: {{timewindow_begin}}
to: {{timewindow_end}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 1
jsonpath "$.records[0].fields.level" == "info"
jsonpath "$.records[0].fields.message" == "Hello, World!"

GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Basic {{admin_token}}
[QueryStringParams]
from: {{timewindow_begin}}
to: {{timewindow_end}}
[Options]
delay: 2s

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 0

DELETE http://localhost:5080/api/v1/streams/default
Authorization: Basic {{admin_token}}
HTTP 200
