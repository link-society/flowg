GET http://localhost:5080/api/v1/transformers
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.transformers" count == 0

GET http://localhost:5080/api/v1/transformers/logfmt
Authorization: Bearer {{admin_token}}
HTTP 404

PUT http://localhost:5080/api/v1/transformers/logfmt
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "script": ". = parse_logfmt!(.message)"
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/transformers/logfmt
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/transformers/logfmt/test
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "level=info message=\"foo bar\""
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.record.level" == "info"
jsonpath "$.record.message" == "foo bar"

DELETE http://localhost:5080/api/v1/transformers/logfmt
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/transformers/logfmt
Authorization: Bearer {{admin_token}}
HTTP 404
