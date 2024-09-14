GET http://localhost:5080/api/v1/streams
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams" isEmpty

PUT http://localhost:5080/api/v1/streams/test
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "config": {
    "indexed_fields": ["level"],
    "ttl": 3600,
    "size": 1024
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/streams
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams.test" exists
jsonpath "$.streams.test.indexed_fields" count == 1
jsonpath "$.streams.test.indexed_fields[0]" == "level"
jsonpath "$.streams.test.ttl" == 3600
jsonpath "$.streams.test.size" == 1024

GET http://localhost:5080/api/v1/streams/test
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.config.indexed_fields" count == 1
jsonpath "$.config.indexed_fields[0]" == "level"
jsonpath "$.config.ttl" == 3600
jsonpath "$.config.size" == 1024

GET http://localhost:5080/api/v1/streams/test_getorcreate
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.config.indexed_fields" count == 0
jsonpath "$.config.ttl" == 0
jsonpath "$.config.size" == 0

GET http://localhost:5080/api/v1/streams
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams.test" exists
jsonpath "$.streams.test_getorcreate" exists

DELETE http://localhost:5080/api/v1/streams/test_getorcreate
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/streams
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams.test" exists
jsonpath "$.streams.test_getorcreate" not exists

DELETE http://localhost:5080/api/v1/streams/test
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/streams
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams" isEmpty
