###############################################################################
# DESCRIPTION: Test stream management through API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that no streams exist
#  - Create a stream "test" with indexed field "level", ttl 3600, and size 1024
#  - Verify the stream appears in the list of streams
#  - Verify the stream configuration

GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams" isEmpty

PUT http://localhost:5080/api/v1/streams/test
Authorization: Bearer {{admin_token}}
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
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams.test" exists
jsonpath "$.streams.test.indexed_fields" count == 1
jsonpath "$.streams.test.indexed_fields[0]" == "level"
jsonpath "$.streams.test.ttl" == 3600
jsonpath "$.streams.test.size" == 1024

GET http://localhost:5080/api/v1/streams/test
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.config.indexed_fields" count == 1
jsonpath "$.config.indexed_fields[0]" == "level"
jsonpath "$.config.ttl" == 3600
jsonpath "$.config.size" == 1024

# -----------------------------------------------------------------------------
# TEST:
#  - Verify that a stream is automatically created if it does not exist
#  - Verify the stream configuration is default (no indexed fields, ttl 0, size 0)
#  - Verify the stream appears in the list of streams
#  - Delete the stream
#  - Verify the stream is no longer in the list of streams

GET http://localhost:5080/api/v1/streams/test_getorcreate
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.config.indexed_fields" count == 0
jsonpath "$.config.ttl" == 0
jsonpath "$.config.size" == 0

GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams.test" exists
jsonpath "$.streams.test_getorcreate" exists

DELETE http://localhost:5080/api/v1/streams/test_getorcreate
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams.test" exists
jsonpath "$.streams.test_getorcreate" not exists

# -----------------------------------------------------------------------------
# TEST:
#  - Delete the "test" stream
#  - Verify the stream is no longer in the list of streams

DELETE http://localhost:5080/api/v1/streams/test
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/streams
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.streams" isEmpty
