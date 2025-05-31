PUT http://localhost:5080/api/v1/forwarders/splunk-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "forwarder": {
    "config": {
      "type": "splunk",
      "endpoint": "http://test-flowg-mockserver:1080/services/collector/event",
      "token": "12345678-1234-1234-1234-123456789012"
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/splunk-success
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "record": {
    "message": "Splunk success test"
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

DELETE http://localhost:5080/api/v1/forwarders/splunk-success
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

PUT http://localhost:5080/api/v1/forwarders/splunk-fail
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "forwarder": {
    "config": {
      "type": "splunk",
      "endpoint": "http://test-flowg-mockserver:1080/services/collector/fail",
      "token": "12345678-1234-1234-1234-123456789012"
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
    "message": "Splunk failure test"
  }
}
HTTP 500

DELETE http://localhost:5080/api/v1/forwarders/splunk-fail
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
