GET http://localhost:9113/cluster/nodes
HTTP 200
[Asserts]
jsonpath "$.nodes" count == 3

GET http://localhost:9114/cluster/nodes
HTTP 200
[Asserts]
jsonpath "$.nodes" count == 3

GET http://localhost:9115/cluster/nodes
HTTP 200
[Asserts]
jsonpath "$.nodes" count == 3
