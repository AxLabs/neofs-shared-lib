package object

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
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

/*
----Object----
+Put
+Get
+Delete
+Head
-Search
-GetRange
-GetRangeHash
*/

func CreateObject(neofsClient *client.NeoFSClient, containerID cid.ID, sessionSigner ecdsa.PrivateKey,
	attributes [][2]string, payload io.Reader) *response.StringResponse {

	ctx := context.Background()

	// region open session

	var prmSession neofsclient.PrmSessionCreate
	// send request to open the session for object writing
	const expirationSession = math.MaxUint64
	prmSession.SetExp(expirationSession)
	prmSession.UseKey(sessionSigner)

	client := neofsClient.LockAndGet()
	resSession, err := client.SessionCreate(ctx, prmSession)
	neofsClient.Unlock()
	if err != nil {
		return response.StringError(err)
	}

	// decode session ID
	var idSession uuid.UUID

	err = idSession.UnmarshalBinary(resSession.ID())
	if err != nil {
		return response.StringError(err)
	}

	// decode session public key
	var keySession neofsecdsa.PublicKey

	err = keySession.Decode(resSession.PublicKey())
	if err != nil {
		return response.StringError(err)
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
		return response.StringError(err)
	}

	// endregion open session

	// pre: tokenSession, signer, context
	var prmPutInit neofsclient.PrmObjectPutInit
	prmPutInit.WithinSession(tokenSession)
	prmPutInit.UseKey(sessionSigner)

	streamObj, err := client.ObjectPutInit(ctx, prmPutInit)
	if err != nil {
		return response.StringError(err)
	}

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

	res, err := streamObj.Close()
	if err != nil {
		return response.StringError(err)
	}
	objectID := res.StoredObjectID()
	return response.NewString(reflect.TypeOf(oid.ID{}), objectID.EncodeToString())
}

func ReadObject(neofsClient *client.NeoFSClient, containerID cid.ID, objectID oid.ID,
	signer ecdsa.PrivateKey) *response.PointerResponse {

	ctx := context.Background()

	var prmGet neofsclient.PrmObjectGet
	prmGet.FromContainer(containerID)
	prmGet.ByID(objectID)
	//prmGet.UseKey(signerDefault)
	prmGet.UseKey(signer)

	client := neofsClient.LockAndGet()
	streamObj, err := client.ObjectGetInit(ctx, prmGet)
	neofsClient.Unlock()
	if err != nil {
		return response.Error(err)
	}
	var b bytes.Buffer
	_ = io.Writer(&b)
	if streamObj.ReadHeader(new(object.Object)) {
		_, err = io.Copy(&b, streamObj)
		if err != nil {
			return response.Error(err)
		}
	}

	_, err = streamObj.Close()
	if err != nil {
		return response.Error(err)
	}
	return response.New(reflect.TypeOf(object.Object{}), b.Bytes())
}

func GetObjectHead(neofsClient *client.NeoFSClient, containerID cid.ID, objectID oid.ID,
	signer ecdsa.PrivateKey) *response.PointerResponse {

	ctx := context.Background()

	var prmObjectHead neofsclient.PrmObjectHead
	prmObjectHead.FromContainer(containerID)
	prmObjectHead.ByID(objectID)
	//prmGet.UseKey(signerDefault)
	prmObjectHead.UseKey(signer)

	client := neofsClient.LockAndGet()
	resObjectHead, err := client.ObjectHead(ctx, prmObjectHead)
	neofsClient.Unlock()
	if err != nil {
		return response.Error(err)
	}

	var objectHeader object.Object
	read := resObjectHead.ReadHeader(&objectHeader)
	if !read {
		return response.Error(fmt.Errorf("could not read object header"))
	}

	v2 := objectHeader.ToV2()
	bytes := v2.StableMarshal(nil)
	return response.New(reflect.TypeOf(object.Object{}), bytes)
}

// DeleteObject marks an object for deletion from the container using NeoFS API protocol.
// As a marker, a special unit called a tombstone is placed in the container.
// It confirms the user's intent to delete the object, and is itself a container object.
// Explicit deletion is done asynchronously, and is generally not guaranteed.
func DeleteObject(neofsClient *client.NeoFSClient, containerID cid.ID, objectID oid.ID,
	signer ecdsa.PrivateKey) *response.StringResponse {

	ctx := context.Background()

	var prmDelete neofsclient.PrmObjectDelete
	prmDelete.FromContainer(containerID)
	prmDelete.ByID(objectID)
	//prmDelete.UseKey(signerDefault)
	prmDelete.UseKey(signer)

	client := neofsClient.LockAndGet()
	res, err := client.ObjectDelete(ctx, prmDelete)
	neofsClient.Unlock()
	if err != nil {
		return response.StringError(err)
	}

	res.Status()
	if !apistatus.IsSuccessful(res.Status()) {
		return response.StringStatusResponse()
	}
	tombstoneID := res.Tombstone()
	return response.NewString(reflect.TypeOf(oid.ID{}), tombstoneID.EncodeToString())
}
