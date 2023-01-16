# Exported Methods

All exported methods from the shared library are described here. You can find specific information
about all exported methods below (e.g., description, parameters, and return type).

### Overview of all exported methods

- [Create client](#create-a-client)
- [Delete client](#delete-a-client)
- [Get account balance](#get-account-balance)
- [Get endpoint information](#get-endpoint-information)
- [Get network information](#get-network-information)
- [Create container](#create-a-container)
- [Get container](#get-a-container)
- [Delete container](#delete-a-container)
- [List containers](#list-owned-containers)
- [Create object](#create-an-object-without-attributes)
- [Create object with 1 attribute](#create-an-object-with-1-attribute)
- [Read object](#read-an-object)
- [Delete object](#delete-an-object)

The structure of the responses in C is defined in the header file [`response.h`](./response.h). The
type `response`
reflects a response that holds a simple string value. In order to be able to parse it correctly, the
type of the return value is returned besides the string. The type `responsePointer` represents a
byte array return value. In order to parse this, it's size and type is provided besides the value.

<hr>

## Client

### Create a client

`CreateClient(privateKey *C.char, neofsEndpoint *C.char) C.pointerResponse`

For all interactions with NeoFS, a client that connects to a GRPC endpoint of a NeoFS node is
needed. For this, you can use the method `CreateClient()` with an endpoint to a NeoFS node and a
private key that will then be used to sign outgoing requests.

#### Parameters

- `privateKey`: hexadecimal string of a private key's big-endian bytes. For example,
  `57dfb170c7cb66b6361ae66021afc1a9a7c5be47aaeb72a82ad5ac963b6da2f9` for the
  WIF `KzAXTwrj1VxQA746zSSMCt9g3omSDfyKnwsayEducuHvKd1LR9mx`.
- `neofsEndpoint`: the GRPC endpoint of a NeoFS storage node. For
  example, `grpcs://st1.t5.fs.neo.org:8082`.

#### Returns

This method returns a uuid (`clientID` in the following) that can be used to connect to
the created NeoFS-Go client in memory. When passing the `clientID` to a method, that corresponding
client is loaded from memory and used to send the corresponding request to the specified endpoint.

### Delete a client

`DeleteClient(clientID *C.char) C.pointerResponse`

#### Parameters

- `clientID`: The client ID.

#### Returns

`1`, if the client was deleted successfully. Otherwise, an error.

<hr>

## Accounting

### Get account balance

`GetBalance(clientID *C.char, publicKey *C.char) C.pointerResponse`

This method gets the NeoFS balance the provided account has on the NeoFS smart contract on the Neo
blockchain.

#### Parameters

- `clientID`: The client ID.
- `publicKey`: encoded compressed public key as hexadecimal. For example,
  `0239e0884856c41605cd1bb660a72ccfc39df8acf576f15a3e593acaa27351b457`.

#### Returns

The return value is based on the type `neo.fs.v2.accounting.Decimal` of
the [NeoFS API](https://github.com/nspcc-dev/neofs-api).

<hr>

## Netmap

### Get endpoint information

`GetEndpoint(clientID *C.char) C.pointerResponse`

Gets information about the endpoint the client is interfacing with. Includes node information and
what version is used in the node.

#### Parameters

- `clientID`: The client ID.

#### Returns

The return value is based on a JSON serialized type `EndpointResponse` (see schema in SCHEMAS.md).

### Get network information

`GetNetworkInfo(clientID *C.char) C.pointerResponse`

Gets the network information the client is interfacing with.

#### Parameters

- `clientID`: The client ID.

#### Returns

The return value is based on the type `neo.fs.v2.netmap.NetworkInfo` of
the [NeoFS API](https://github.com/nspcc-dev/neofs-api).

<hr>

## Container

### Create a container

`PutContainer(clientID *C.char, v2Container *C.char) C.response`

Creates a container.

#### Parameters

- `clientID`: The client ID.
- `v2Container`: The container to create. This parameter has to be based on the
  type `neo.fs.v2.Container` of the [NeoFS API](https://github.com/nspcc-dev/neofs-api) and is
  required to hold an `OwnerID`.

#### Returns

The container ID.

### Get a container

`GetContainer(clientID *C.char, containerID *C.char) C.pointerResponse`

Gets the container header.

#### Parameters

- `clientID`: The client ID.
- `containerID`: The container ID.

#### Returns

The return value is based on the type `neo.fs.v2.container.Container` of
the [NeoFS API](https://github.com/nspcc-dev/neofs-api).

### Delete a container

`DeleteContainer(clientID *C.char, containerID *C.char) C.pointerResponse`

Deletes a container.

#### Parameters

- `clientID`: The client ID.
- `containerID`: The container ID.

#### Returns

`1`, if the container was deleted successfully. Otherwise, an error.

### List owned containers

`ListContainer(clientID *C.char, ownerPubKey *C.char) C.pointerResponse`

Returns a list of containers owned by the provided public key (`ownerPubKey`).

#### Parameters

- `clientID`: The client ID.
- `ownerPubKey`: The owner's encoded compressed public key as hexadecimal. For example,
  `0239e0884856c41605cd1bb660a72ccfc39df8acf576f15a3e593acaa27351b457`.

#### Returns

A byte array that consists of concatenated container IDs. One container ID is 33 bytes long.

<hr>

## Object

### Create an object without attributes

`CreateObjectWithoutAttributes(clientID *C.char, containerID *C.char, fileBytes unsafe.Pointer, fileSize C.int, sessionSignerPrivKey *C.char) C.response`

Creates an object within the provided container. Does not add any attributes to the object.

#### Parameters

- `clientID`: The client ID.
- `containerID`: The container ID.
- `fileBytes`: The bytes to store as an object in the container.
- `fileSize`: The number of bytes of `fileBytes`.
- `sessionSignerPrivKey`: The signer that opens the session to write the object.

#### Returns

The object ID.

### Create an object with 1 attribute

`CreateObject(clientID *C.char, containerID *C.char, fileBytes unsafe.Pointer, fileSize C.int, sessionSignerPrivKey *C.char, attributeKey *C.char, attributeValue *C.char) C.response`

Creates an object within the provided container. Adds exactly one attribute to the object.

> Note: Adding multiple attributes to an object is currently not yet supported.

#### Parameters

- `clientID`: The client ID.
- `containerID`: The container ID.
- `fileBytes`: The bytes to store as an object in the container.
- `fileSize`: The number of bytes of `fileBytes`.
- `sessionSignerPrivKey`: The signer that opens the session to write the object.
- `attributeKey`: The attribute key.
- `attributeValue`: The attribute value.

#### Returns

The object ID.

### Read an object

`ReadObject(clientID *C.char, containerID *C.char, objectID *C.char, signer *C.char) C.pointerResponse`

Reads the object with objectID from the container with containerID.

#### Parameters

- `clientID`: The client ID.
- `containerID`: The container ID.
- `objectID`: The object ID.
- `signer`: The signer that opens the session to read the object.

#### Returns

The object bytes.

### Delete an object

`DeleteObject(clientID *C.char, containerID *C.char, objectID *C.char, signer *C.char) C.response`

Deletes an object with objectID from the container with containerID.

#### Parameters

- `clientID`: The client ID.
- `containerID`: The container ID.
- `objectID`: The object ID.
- `signer`: The signer that opens the session to read the object.

#### Returns

The tombstone object ID.

> When deleting an object, an associated `tombstone` object is created to mark an object as deleted
> for potential replicas in the network. After some time, that tombstone will be replicated on nodes
> that are storing the object and eventually the nodes will delete both the initial object and the
> associated tombstone. In some way, it's like a pending delete request.
