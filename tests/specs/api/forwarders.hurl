GET http://localhost:5080/api/v1/forwarders
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarders" count == 0

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
HTTP 404

PUT http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
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
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarders" count == 1
jsonpath "$.forwarders[0]" == "httpbin"

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.forwarder.config.type" == "http"
jsonpath "$.forwarder.config.url" == "http://httpbin.org/anything"
jsonpath "$.forwarder.config.headers.Foo" == "Bar"

DELETE http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/forwarders/httpbin
Authorization: Bearer {{admin_token}}
HTTP 404
