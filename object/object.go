package object

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	"io"
	"math"
	"reflect"

	"github.com/AxLabs/neofs-api-shared-lib/client"
	"github.com/google/uuid"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	cid "github.com/nspcc-dev/neofs-sdk-go/container/id"
	neofsecdsa "github.com/nspcc-dev/neofs-sdk-go/crypto/ecdsa"
	"github.com/nspcc-dev/neofs-sdk-go/object"
	oid "github.com/nspcc-dev/neofs-sdk-go/object/id"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	"github.com/nspcc-dev/neofs-sdk-go/user"
)

func CreateObject(neofsClient *client.NeoFSClient, containerID cid.ID, sessionSigner ecdsa.PrivateKey,
	attributes [][2]string, payload io.Reader) *response.StringResponse {

	// How do we pass more attributes to C (dynamic number of attributes)?

	ctx := context.Background()

	// region open session

	var prmSession neofsclient.PrmSessionCreate
	// send request to open the session for object writing
	const expirationSession = math.MaxUint64
	prmSession.SetExp(expirationSession)
	prmSession.UseKey(sessionSigner)

	client := neofsClient.LockAndGet()
	resSession, err := client.SessionCreate(ctx, prmSession)
	if err != nil {
		return response.StringError(err)
		//fmt.Errorf("open session with the remote node: %w", err)
	}
	neofsClient.Unlock()

	// decode session ID
	var idSession uuid.UUID

	err = idSession.UnmarshalBinary(resSession.ID())
	if err != nil {
		return response.StringError(err)
		//fmt.Errorf("invalid session ID in session response: %w", err)
	}

	// decode session public key
	var keySession neofsecdsa.PublicKey

	err = keySession.Decode(resSession.PublicKey())
	if err != nil {
		return response.StringError(err)
		//fmt.Errorf("invalid session public key in session response: %w", err)
	}

	// form token of the object session
	var tokenSession session.Object
	tokenSession.SetID(idSession)
	tokenSession.SetExp(expirationSession)
	tokenSession.BindContainer(containerID) // prm.Container
	tokenSession.ForVerb(session.VerbObjectPut)
	tokenSession.SetAuthKey(&keySession)

	// sign the session token
	err = tokenSession.Sign(sessionSigner)
	if err != nil {
		return response.StringError(err) //fmt.Errorf("sign session token: %w", err)
	}

	fmt.Println("session created")

	// endregion open session

	// pre: tokenSession, signer, context
	var prmPutInit neofsclient.PrmObjectPutInit
	prmPutInit.WithinSession(tokenSession)
	prmPutInit.UseKey(sessionSigner)

	streamObj, err := client.ObjectPutInit(ctx, prmPutInit)
	if err != nil {
		return response.StringError(err)
	}

	fmt.Println("object put initialized")

	var idCreator user.ID
	user.IDFromKey(&idCreator, sessionSigner.PublicKey)
	var obj object.Object
	obj.SetContainerID(containerID)
	obj.SetOwnerID(&idCreator)

	// add attributes
	if attributes != nil {
		attrs := make([]object.Attribute, len(attributes))

		for i := range attributes {
			attrs[i].SetKey(attributes[i][0])
			attrs[i].SetValue(attributes[i][1])
		}

		obj.SetAttributes(attrs...)
	}

	v, set := obj.ID()
	if !set {
		fmt.Println("no id set in object header")
	} else {
		fmt.Println("id set in object header")
	}
	fmt.Println("object id: " + v.String())

	fmt.Println("v2 object id: " + string(obj.ToV2().GetObjectID().GetValue()))

	if streamObj.WriteHeader(obj) && payload != nil {
		var n int
		buf := make([]byte, 100<<10)
		for {
			n, err = payload.Read(buf)
			if n > 0 {
				if !streamObj.WritePayloadChunk(buf[:n]) {
					break
				}
				continue
			}
			if errors.Is(err, io.EOF) {
				break
			}
			return response.StringError(err) // read payload
		}
	}

	fmt.Println("stream obj written")

	res, err := streamObj.Close()
	if err != nil {
		return response.StringError(err)
	}
	objectID := res.StoredObjectID()
	fmt.Println(objectID.String())
	return response.NewString(reflect.TypeOf(oid.ID{}), objectID.EncodeToString())
}

func ReadObject(neofsClient *client.NeoFSClient, containerID cid.ID, objectID oid.ID,
	signer ecdsa.PrivateKey) *response.PointerResponse {

	client := neofsClient.LockAndGet()
	ctx := context.Background()

	var prmGet neofsclient.PrmObjectGet
	prmGet.FromContainer(containerID)
	prmGet.ByID(objectID)
	//prmGet.UseKey(signerDefault)
	prmGet.UseKey(signer)

	fmt.Println("prm set to read")

	streamObj, err := client.ObjectGetInit(ctx, prmGet)
	if err != nil {
		return response.Error(err)
	}
	fmt.Println("streamObj initialized")
	var b bytes.Buffer
	_ = io.Writer(&b)
	//foo := bufio.NewWriter(&b)
	if streamObj.ReadHeader(new(object.Object)) {
		fmt.Println("read header success")
		_, err = io.Copy(&b, streamObj)
		fmt.Println("io copy success")
		if err != nil {
			return response.Error(err)
		}
	}
	fmt.Println("header written")

	_, err = streamObj.Close()
	if err != nil {
		return response.Error(err)
	}
	return response.New(reflect.TypeOf(object.Object{}), b.Bytes())
}

//import (
//	"context"
//	"fmt"
//	"github.com/AxLabs/neofs-api-shared-lib/client"
//	"github.com/AxLabs/neofs-api-shared-lib/response"
//	"github.com/google/uuid"
//	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
//	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
//	"github.com/nspcc-dev/neofs-sdk-go/object"
//	oid "github.com/nspcc-dev/neofs-sdk-go/object/id"
//	"reflect"
//	"sync"
//	"unsafe"
//)
//
///*
//----Object----
//Put
//Get
//Delete
//Head
//Search
//GetRange
//GetRangeHash
//*/
//
////export ObjectPutInit
//func ObjectPutInit(clientID *C.char) C.pointerResponse {
//	ctx := context.Background()
//
//	var prmObjectPutInit neofsclient.PrmObjectPutInit
//	// 	prmObjectPutInit.SetCopiesNumber(uint32)
//
//	neofsClient, err := client.GetClient(clientID)
//	if err != nil {
//		return response.PointerResponseClientError()
//	}
//  client := neofsClient.LockAndGet()
//	w, err := client.ObjectPutInit(ctx, prmObjectPutInit)
//	neofsClient.Unlock()
//	if err != nil {
//		return response.PointerResponseError(err.Error())
//	}
//
//	u, err := uuid.NewUUID()
//	if err != nil {
//		return response.PointerResponseError(err.Error())
//	}
//
//	//w.MarkLocal() // optional
//	//w.UseKey() // optional
//
//	if wMap == nil {
//		initWriterMap(u, w)
//	} else {
//		wMap.put(u, w)
//	}
//	return response.PointerResponse(reflect.TypeOf(u), []byte(u.String()))
//}
//
////export WriteHeader
//func WriteHeader(writerID *C.char, header *C.char) C.pointerResponse {
//	w, err := getWriter(writerID)
//	if err != nil {
//		return response.PointerResponseError("could not get writer")
//	}
//	w.mu.Lock()
//	//w.writer.WithXHeaders()
//	//w.writer.WithinSession()
//	//w.writer.WithBearerToken()
//	o := object.New()
//	written := w.writer.WriteHeader(*o)
//	w.mu.Unlock()
//
//	return response.PointerResponseBoolean(written)
//}
//
//// Can be called multiple times (appended). If the chunk exceeds the maximal size, it is
//// automatically split (assumption in neofs-sdk-go: maximal chunk size of 3MB).
////export WritePayloadChunk
//func WritePayloadChunk(writerID *C.char, chunk unsafe.Pointer, chunkSize C.int) C.pointerResponse {
//	bytes := C.GoBytes(chunk, chunkSize)
//	fmt.Println(bytes)
//
//	w, err := getWriter(writerID)
//	if err != nil {
//		return response.PointerResponseError("could not get writer")
//	}
//	w.mu.Lock()
//	//w.writer.WithXHeaders()
//	//w.writer.WithinSession()
//	//w.writer.WithBearerToken()
//	written := w.writer.WritePayloadChunk(bytes)
//	w.mu.Unlock()
//
//	return response.PointerResponseBoolean(written)
//}
//
////export CloseWriter
//func CloseWriter(writerID *C.char) C.response {
//
//	w, err := getWriter(writerID)
//	if err != nil {
//		return response.ResponseError("could not get writer")
//	}
//	w.mu.Lock()
//	//w.writer.WithXHeaders()
//	//w.writer.WithinSession()
//	//w.writer.WithBearerToken()
//	res, err := w.writer.Close()
//	if err != nil {
//		return response.ResponseError(err.Error())
//	}
//	wMap.delete(writerID)
//	w.mu.Unlock()
//
//	if !apistatus.IsSuccessful(res.Status()) {
//		return response.ResultStatusErrorResponse()
//	}
//
//	var oid oid.ID
//	read := res.ReadStoredObjectID(&oid)
//	if !read {
//		return response.ResponseError("could not read stored object id")
//	}
//	return response.Response(reflect.TypeOf(oid), oid.String())
//}
//
////region client
//
//var wMap *WriterMap
//
//type Writer struct {
//	mu     sync.RWMutex
//	writer *neofsclient.ObjectWriter
//}
//
//type WriterMap struct {
//	mu      sync.RWMutex
//	writers map[uuid.UUID]*Writer
//}
//
//func (wMap *WriterMap) put(id uuid.UUID, newWriter *neofsclient.ObjectWriter) {
//	wMap.mu.Lock()
//	wMap.writers[id] = &Writer{
//		mu:     sync.RWMutex{},
//		writer: newWriter,
//	}
//	wMap.mu.Unlock()
//}
//
//func (wMap *WriterMap) delete(writerID *C.char) bool {
//	id, err := uuid.Parse(C.GoString(writerID))
//	if err != nil {
//		return false
//	}
//	wMap.mu.Lock()
//	delete(wMap.writers, id)
//	wMap.mu.Unlock()
//	return true
//}
//
//func initWriterMap(id uuid.UUID, newWriter *neofsclient.ObjectWriter) {
//	wMap = &WriterMap{sync.RWMutex{}, map[uuid.UUID]*Writer{id: {sync.RWMutex{}, newWriter}}}
//}
//
//func getWriter(writerID *C.char) (*Writer, error) {
//	if wMap == nil {
//		return nil, fmt.Errorf("no writer present")
//	}
//	wID, err := uuid.Parse(C.GoString(writerID))
//	if err != nil {
//		return nil, fmt.Errorf("could not parse provided writer id. id was " + C.GoString(writerID))
//	}
//	wMap.mu.Lock()
//	w := wMap.writers[wID]
//	if w == nil {
//		return nil, fmt.Errorf("no client present with id %v", C.GoString(writerID))
//	}
//	wMap.mu.Unlock()
//	return w, nil
//}
//
//////export GetObjectInit
////func GetObjectInit(clientID *C.char, v2ContainerID *C.char) *C.PointerResponse {
////	cli, err := GetClient(clientID)
////	if err != nil {
////		return PointerResponseClientError()
////	}
////	cli.mu.RLock()
////	ctx := context.Background()
////
////	// Parse the container
////	id, err := getContainerIDFromV2(v2ContainerID)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////
////	var prmObjectGet neofsclient.PrmObjectGet
////	prmObjectGet.FromContainer(*id) // required
////	prmObjectGet.ByID()             // required
////	//prmObjectGet.MarkLocal()        // optional, tells the server to execute operation locally
////	//prmObjectGet.MarkRaw()          // optional, marks intent to read physically stored object
////	//prmObjectGet.WithinSession()    // optional
////	//prmObjectGet.WithBearerToken()  // optional
////
////	Response, err := cli.client.ObjectGetInit(ctx, prmObjectGet)
////	cli.mu.RUnlock()
////	if err != nil {
////		panic(err)
////	}
////
////	// todo: Check how exactly to read object bytes
////	//Response.UseKey()
////	//Response.Read()
////	//Response.ReadChunk()
////	//Response.ReadHeader()
////	read, err := Response.Read()
////	if err != nil {
////		return nil
////	}
////	return C.CString() // return pointer to object reader and
////}
//
//////export DeleteObject
////func DeleteObject(clientID *C.char, v2ContainerID *C.char, v2ObjectID *C.char, v2SessionToken *C.char,
////	v2BearerToken *C.char) C.PointerResponse {
////
////	cli, err := GetClient(clientID)
////	if err != nil {
////		return PointerResponseClientError()
////	}
////	cli.mu.RLock()
////	ctx := context.Background()
////	containerID, err := getContainerIDFromV2(v2ContainerID)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	objectID, err := getObjectIDFromV2(v2ObjectID)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	sessionToken, err := getSessionTokenFromV2(v2SessionToken)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	bearerToken, err := getBearerTokenFromV2(v2BearerToken)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	var prmObjectDelete neofsclient.PrmObjectDelete
////	prmObjectDelete.FromContainer(*containerID)
////	prmObjectDelete.ByID(*objectID)
////	prmObjectDelete.WithinSession(*sessionToken)
////	prmObjectDelete.WithBearerToken(*bearerToken)
////
////	resObjectDelete, err := cli.client.ObjectDelete(ctx, prmObjectDelete)
////	cli.mu.RUnlock()
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	if !apistatus.IsSuccessful(resObjectDelete.Status()) {
////		return ResultStatusErrorResponsePointer()
////	}
////	readTombStoneID := new(oid.ID)
////	tombstoneRead := resObjectDelete.ReadTombstoneID(readTombStoneID)
////	if !tombstoneRead {
////		return PointerResponseError("could not read object's tombstone")
////	}
////	json, err := readTombStoneID.MarshalJSON()
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	return PointerResponse(reflect.TypeOf(tombstoneRead), json)
////}
//
//////export GetObjectHead
////func GetObjectHead(clientID *C.char, v2ContainerID *C.char, v2ObjectID *C.char, v2SessionToken *C.char,
////	v2BearerToken *C.char) C.PointerResponse {
////
////	cli, err := GetClient(clientID)
////	if err != nil {
////		return PointerResponseClientError()
////	}
////	cli.mu.RLock()
////	ctx := context.Background()
////	containerID, err := getContainerIDFromV2(v2ContainerID)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	objectID, err := getObjectIDFromV2(v2ObjectID)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	sessionToken, err := getSessionTokenFromV2(v2SessionToken)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	bearerToken, err := getBearerTokenFromV2(v2BearerToken)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	var prmObjectHead neofsclient.PrmObjectHead
////	prmObjectHead.FromContainer(*containerID)
////	prmObjectHead.ByID(*objectID)
////	prmObjectHead.WithinSession(*sessionToken)
////	prmObjectHead.WithBearerToken(*bearerToken)
////
////	resObjectHead, err := cli.client.ObjectHead(ctx, prmObjectHead)
////	cli.mu.RUnlock()
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	if !apistatus.IsSuccessful(resObjectHead.Status()) {
////		return ResultStatusErrorResponsePointer()
////	}
////	dst := new(object.Object)
////	resObjectHead.ReadHeader(dst)
////	json, err := dst.MarshalJSON()
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	return PointerResponse(reflect.TypeOf(dst), json)
////}
//
//////export SearchObject
////func SearchObject(clientID *C.char, v2ContainerID *C.char, v2SessionToken *C.char, v2BearerToken *C.char, v2Filters *C.char) C.Response {
////	cli, err := GetClient(clientID)
////	if err != nil {
////		return ResponseClientError()
////	}
////	cli.mu.RLock()
////	ctx := context.Background()
////
////	containerID, err := getContainerIDFromV2(v2ContainerID)
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////	sessionToken, err := getSessionTokenFromV2(v2SessionToken)
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////	bearerToken, err := getBearerTokenFromV2(v2BearerToken)
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////	filters, err := getFiltersFromV2(v2Filters)
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////	var prmObjectSearch neofsclient.PrmObjectSearch
////	prmObjectSearch.InContainer(*containerID)
////	prmObjectSearch.WithinSession(*sessionToken)
////	prmObjectSearch.WithBearerToken(*bearerToken)
////	prmObjectSearch.SetFilters(*filters)
////	//prmObjectSearch.MarkLocal()
////
////	resObjectSearchInit, err := cli.client.ObjectSearchInit(ctx, prmObjectSearch)
////	cli.mu.RUnlock()
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////
////	//resObjectSearchInit.UseKey()
////	//resObjectSearchInit.Read()
////	//resObjectSearchInit.Close()
////
////	read, b := resObjectSearchInit.Read()
////	return Response("SearchObject", read)
////}
//
//////export GetRange
////func GetRange(clientID *C.char, v2ContainerID *C.char, v2ObjectID *C.char, v2SessionToken *C.char,
////	v2BearerToken *C.char, offset *C.char, length *C.char) C.Response {
////
////	cli, err := GetClient(clientID)
////	if err != nil {
////		return ResponseClientError()
////	}
////	cli.mu.RLock()
////	ctx := context.Background()
////
////	containerID, err := getContainerIDFromV2(v2ContainerID)
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////	objectID, err := getObjectIDFromV2(v2ObjectID)
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////	sessionToken, err := getSessionTokenFromV2(v2SessionToken)
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////	bearerToken, err := getBearerTokenFromV2(v2BearerToken)
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////
////	var prmObjectRange neofsclient.PrmObjectRange
////	prmObjectRange.FromContainer(*containerID)
////	prmObjectRange.ByID(*objectID)
////	prmObjectRange.WithinSession(*sessionToken)
////	prmObjectRange.WithBearerToken(*bearerToken)
////	prmObjectRange.SetLength(length)
////	prmObjectRange.SetOffset(offset)
////
////	Response, err := cli.client.ObjectRangeInit(ctx, prmObjectRange)
////	cli.mu.RUnlock()
////	if err != nil {
////		return ResponseError(err.Error())
////	}
////
////	Response.Read()
////	return Response("GetRange")
////}
//
//////export GetRangeHash
////func GetRangeHash(clientID *C.char, v2ContainerID *C.char) {
////}
