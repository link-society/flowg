GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer pat:{{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarders" count == 0

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer pat:{{admin_token}}
HTTP 404

PUT http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer pat:{{admin_token}}
Content-Type: application/json
{
  "forwarder": {
    "config": {
      "type": "http",
      "url": "http://httpbin.org/anything",
      "headers": {
        "Foo": "Bar"
      }
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer pat:{{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarders" count == 1
jsonpath "$.forwarders[0]" == "httpbin"

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer pat:{{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarder.config.type" == "http"
jsonpath "$.forwarder.config.url" == "http://httpbin.org/anything"
jsonpath "$.forwarder.config.headers.Foo" == "Bar"

POST http://localhost:5080/api/v1/test/forwarders/httpbin
Authorization: Bearer pat:{{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "Hello, World!"
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

DELETE http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer pat:{{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer pat:{{admin_token}}
HTTP 404

PUT http://localhost:5080/api/v1/forwarders/fail
Authorization: Bearer pat:{{admin_token}}
Content-Type: application/json
{
  "forwarder": {
    "config": {
      "type": "http",
      "url": "http://httpbin.org/status/500",
      "headers": {
        "Foo": "Bar"
      }
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/forwarders/fail
Authorization: Bearer pat:{{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "Hello, World!"
  }
}
HTTP 500

DELETE http://localhost:5080/api/v1/forwarders/fail
Authorization: Bearer pat:{{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
