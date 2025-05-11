POST http://localhost:5080/api/v1/pipelines/default/logs/struct
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "records": [
    {
      "level": "info",
      "message": "test"
    }
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 1

# Node 0
GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
[Options]
retry: 10
retry-interval: 200ms

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 1
jsonpath "$.records[0].fields.level" == "info"
jsonpath "$.records[0].fields.message" == "test"

# Node 1
GET http://localhost:5081/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
[Options]
retry: 10
retry-interval: 200ms

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 1
jsonpath "$.records[0].fields.level" == "info"
jsonpath "$.records[0].fields.message" == "test"

# Node 2
GET http://localhost:5082/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}
[Options]
retry: 10
retry-interval: 200ms

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 1
jsonpath "$.records[0].fields.level" == "info"
jsonpath "$.records[0].fields.message" == "test"

# --- Cleanup -----------------------------------------------------------------

DELETE http://localhost:5080/api/v1/streams/default
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Node 0
GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{admin_token}}
[Options]
retry: 10
retry-interval: 200ms

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams" isEmpty

# Node 1
GET http://localhost:5081/api/v1/streams
Authorization: Bearer {{admin_token}}
[Options]
retry: 10
retry-interval: 200ms

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams" isEmpty

# Node 2
GET http://localhost:5082/api/v1/streams
Authorization: Bearer {{admin_token}}
[Options]
retry: 10
retry-interval: 200ms

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams" isEmpty
