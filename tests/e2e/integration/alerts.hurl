GET http://localhost:5080/api/v1/alerts
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.alerts" count == 0

GET http://localhost:5080/api/v1/alerts/httpbin
Authorization: Basic {{admin_token}}
HTTP 404

PUT http://localhost:5080/api/v1/alerts/httpbin
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "webhook": {
    "url": "http://httpbin.org/anything",
    "headers": {
      "Foo": "Bar"
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/alerts
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.alerts" count == 1
jsonpath "$.alerts[0]" == "httpbin"

GET http://localhost:5080/api/v1/alerts/httpbin
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/alerts/httpbin
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "Hello, World!"
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

DELETE http://localhost:5080/api/v1/alerts/httpbin
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true

GET http://localhost:5080/api/v1/alerts/httpbin
Authorization: Basic {{admin_token}}
HTTP 404

PUT http://localhost:5080/api/v1/alerts/fail
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "webhook": {
    "url": "http://httpbin.org/status/500",
    "headers": {
      "Foo": "Bar"
    }
  }
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

POST http://localhost:5080/api/v1/test/alerts/fail
Authorization: Basic {{admin_token}}
Content-Type: application/json
{
  "record": {
    "message": "Hello, World!"
  }
}
HTTP 500

DELETE http://localhost:5080/api/v1/alerts/fail
Authorization: Basic {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
