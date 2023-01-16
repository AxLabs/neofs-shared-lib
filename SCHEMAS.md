# Response Schemas

Some responses from the shared library have a custom structure to match the requirements of
procedural programming paradigm (not object oriented). Currently, this affects the responses
from `GetEndpoint()` and `ListContainers()` which are described below.

### `GetEndpoint()`

```json
{
  "title": "Endpoint Response",
  "description": "",
  "type": "object",
  "netmap.NodeInfo": {
    "description": "The Node Information object based on the Protocol Buffers JSON type `neo.fs.v2.netmap.NodeInfo` of the neofs-api (marshalled to JSON).",
    "type": "object"
  },
  "version.Version": {
    "description": "The Node Version object based on the Protocol Buffers JSON type `neo.fs.v2.refs.Version` of the neofs-api (marshalled to JSON).",
    "type": "object"
  }
}
```

### `ListContainers()`

```json
{
  "title": "ListContainer Response",
  "description": "",
  "type": "object",
  "containers": {
    "description": "List of container IDs of all owned containers",
    "type": "array",
    "items": {
      "description": "Container IDs",
      "type": "string"
    }
  }
}
```
