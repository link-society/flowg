POST http://localhost:5080/api/v1/test/transformer
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "code": ". = parse_logfmt!(.message)",
  "record": {
    "message": "level=info message=\"foo bar\""
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.record.level" == "info"
jsonpath "$.record.message" == "foo bar"

GET http://localhost:5080/api/v1/transformers
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.transformers" count == 0

GET http://localhost:5080/api/v1/transformers/logfmt
Authorization: Basic {{admin_token}}
HTTP 404

PUT http://localhost:5080/api/v1/transformers/logfmt
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "script": ". = parse_logfmt!(.message)"
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/transformers
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.transformers" count == 1
jsonpath "$.transformers[0]" == "logfmt"

GET http://localhost:5080/api/v1/transformers/logfmt
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

DELETE http://localhost:5080/api/v1/transformers/logfmt
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/transformers/logfmt
Authorization: Basic {{admin_token}}
HTTP 404
