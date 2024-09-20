POST http://localhost:5080/api/v1/pipelines/default/logs
Authorization: Bearer pat:{{admin_token}}
Content-Type: application/json
{
  "record": {
    "level": "info",
    "message": "Hello, World!"
  }
}
HTTP 200
