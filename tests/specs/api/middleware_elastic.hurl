###############################################################################
# DESCRIPTION: Test ElasticSearch compatibility API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Check that an index (a FlowG pipeline) exists
#  - Verify that the endpoint returns the proper ElasticSearch headers

HEAD http://localhost:5080/api/v1/middlewares/elastic/default
[BasicAuth]
root: root
HTTP 200
X-Elastic-Product: Elasticsearch

# -----------------------------------------------------------------------------
# TEST:
#  - Send logs to the default pipeline via the ElasticSearch API
#  - Verify logs have been ingested
#  - Purge the stream

POST http://localhost:5080/api/v1/middlewares/elastic/default/_doc
Content-Type: application/json
[BasicAuth]
root: root
```
{
  "foo": {
    "bar": "baz"
  }
}
```
HTTP 200
X-Elastic-Product: Elasticsearch

GET http://localhost:5080/api/v1/streams/default/logs
Authorization: Bearer {{admin_token}}
[Query]
from: {{timewindow_begin}}
to: {{timewindow_end}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.records" count == 1
jsonpath "$.records[*].fields['foo.bar']" contains "baz"

DELETE http://localhost:5080/api/v1/streams/default
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
