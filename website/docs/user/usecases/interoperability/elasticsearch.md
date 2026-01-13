---
sidebar_position: 1
---

# ElasticSearch API

**FlowG** supports a subset of the ElasticSearch API, allowing you to plug it
where you would use [ElasticSearch](https://www.elastic.co/elasticsearch),
without any change in your application.

The compatibility API is available under the follwing endpoint:
`/api/v1/middlewares/elastic`.

## Configure the ElasticSearch client

> **NB:** Adapt the username/password and URL according to your setup.

In Go:

```go
cfg := elasticsearch.Config{
  Username:  "root",
  Password:  "root",
  Addresses: []string{"http://localhost:5080/api/v1/middlewares/elastic/"},
}
client, err := elasticsearch.NewClient(cfg)
```

In Javascript:

```javascript
const client = new Client({
  node: 'http://localhost:5080/api/v1/middlewares/elastic/',
  auth: {
    username: 'root',
    password: 'root',
  },
})
```

In Python:

```python
client = Elasticsearch(
  hosts=["http://localhost:5080/api/v1/middlewares/elastic/"],
  basic_auth=("root", "root"),
)
```

## Supported authentication methods

### HTTP Basic

The given credentials map directly to FlowG users.

### :construction: HTTP Bearer

> :warning: **This feature is not yet supported.** :warning:

The given token maps to either a FlowG JSON Web Token, or a FlowG Personal
Access Token.

## Supported operations

### Check if index exists

https://www.elastic.co/docs/api/doc/elasticsearch/operation/operation-indices-exists

```
HEAD /api/v1/middlewares/elastic/{index}
```

> **NB:** The name of the index maps to the name of a FlowG pipeline

| Response | When |
| --- | --- |
| `401 Unauthorized` | The user could not be authenticated (does not exist, or invalid password) |
| `403 Forbidden` | The user does not have the `read_pipelines` permission |
| `404 Not Found` | The pipeline does not exist |
| `200 0K` | The pipeline exists |
| `500 Internal Server Error` | An error occured in FlowG |

**Example usage:**

In Go:

```go
resp, err := client.Indices.Exists(
  []string{"test"},
  client.Indices.Exists.WithContext(context.TODO())
)
```

In Javascript:

```javascript
resp = client.indices.exists({ index: 'test' })
```

In Python:

```python
resp = client.indices.exists(index="test")
```

### Index document

https://www.elastic.co/docs/api/doc/elasticsearch/operation/operation-index

```
POST /api/v1/middlewares/elastic/{index}/_doc
{
  "@timestamp": "...",
  "message": "..."
}
```

> **NB:** The name of the index maps to the name of a FlowG pipeline

| Response | When |
| --- | --- |
| `401 Unauthorized` | The user could not be authenticated (does not exist, or invalid password) |
| `403 Forbidden` | The user does not have the `send_logs` permission |
| `400 Bad Request` | The request body was not JSON |
| `200 0K` | The document (in the request body) was successfully processed through the pipeline |
| `500 Internal Server Error` | An error occured in FlowG |

> **NB:** Since FlowG's datamodel is flat, the document will be flattenned first:

```json
{
  "foo": {
    "bar": "baz"
  }
}
```

will become

```json
{
  "foo.bar": "baz"
}
```

**Example usage:**

In Go:

```go
resp, err = client.Index(
  "test",
  bytes.NewReader([]byte(`{"message": "hello world"}`)),
  client.Index.WithContext(context.TODO()),
)
```

In Javascript:

```javascript
resp = client.index({
  index: 'test',
  document: {foo: {bar: 'baz'}},
})
```

In Python:

```python
resp = client.index(
  index="test",
  document={"foo": {"bar": "baz"}},
)
```

## Roadmap

You can find the tracking issue on Github
[here](https://github.com/link-society/flowg/issues/853).

Feel free to make new feature requests, or report bugs on the existing supported
operations.
