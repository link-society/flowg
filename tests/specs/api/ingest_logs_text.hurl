POST http://localhost:5080/api/v1/pipelines/default/logs/text
Authorization: Bearer {{admin_token}}
Content-Type: text/plain
```
hello
world
```
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.processed_count" == 2

DELETE http://localhost:5080/api/v1/streams/default
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
