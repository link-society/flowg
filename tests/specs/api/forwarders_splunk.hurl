GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarders" count == 0

PUT http://localhost:5080/api/v1/forwarders/splunk-test
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder": {
    "version": 2,
    "config": {
      "type": "splunk",
      "endpoint": "http://httpbin.org/anything",
      "token": "test-token"
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarders" count == 1
jsonpath "$.forwarders[0]" == "splunk-test"

GET http://localhost:5080/api/v1/forwarders/splunk-test
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarder.config.type" == "splunk"
jsonpath "$.forwarder.config.endpoint" == "http://httpbin.org/anything"
jsonpath "$.forwarder.config.token" == "test-token"

POST http://localhost:5080/api/v1/test/forwarders/splunk-test
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "Hello, World!"
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

DELETE http://localhost:5080/api/v1/forwarders/splunk-test
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/forwarders/splunk-test
Authorization: Bearer {{admin_token}}
HTTP 404

PUT http://localhost:5080/api/v1/forwarders/splunk-fail
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "forwarder": {
    "version": 2,
    "config": {
      "type": "splunk",
      "endpoint": "http://httpbin.org/status/500",
      "token": "test-token"
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/splunk-fail
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "Hello, World!"
  }
}
HTTP 500

DELETE http://localhost:5080/api/v1/forwarders/splunk-fail
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true 