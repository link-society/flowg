###############################################################################
# DESCRIPTION: Test transformer management through API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that it is possible to test a transformer before creating it

POST http://localhost:5080/api/v1/test/transformer
Authorization: Bearer {{admin_token}}
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

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that no transformers exist
#  - Create a transformer "logfmt" with a script that parses logfmt messages
#  - Verify the transformer appears in the list of transformers

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

GET http://localhost:5080/api/v1/transformers
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.transformers" count == 1
jsonpath "$.transformers[0]" == "logfmt"

GET http://localhost:5080/api/v1/transformers/logfmt
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# -----------------------------------------------------------------------------
# TEST:
#  - Delete the "logfmt" transformer
#  - Verify it no longer appears in the list of transformers

DELETE http://localhost:5080/api/v1/transformers/logfmt
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/transformers/logfmt
Authorization: Bearer {{admin_token}}
HTTP 404
