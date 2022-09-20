package main

/*
#include <stdlib.h>

#ifndef RESPONSE_H
#define RESPONSE_H
#include "response.h"
#endif
*/
import "C"
import (
	"context"
	v2container "github.com/nspcc-dev/neofs-api-go/v2/container"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	cid "github.com/nspcc-dev/neofs-sdk-go/container/id"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	"reflect"
)

/*
----Container----
Put
Get
Delete
List
SetExtendedACL
GetExtendedACL
AnnounceUsedSpace
*/

//region container

//export PutContainer
func PutContainer(clientID *C.char, v2Container *C.char) C.response {
	ctx := context.TODO()
	cnr, err := getContainerFromC(v2Container)
	if err != nil {
		return responseError(err.Error())
	}

	var prmContainerPut neofsclient.PrmContainerPut
	prmContainerPut.SetContainer(*cnr)

	neofsClient, err := getClient(clientID)
	if err != nil {
		return responseClientError()
	}
	neofsClient.mu.Lock()
	resContainerPut, err := neofsClient.client.ContainerPut(ctx, prmContainerPut)
	neofsClient.mu.Unlock()
	if err != nil {
		return responseError(err.Error())
	}

	if !apistatus.IsSuccessful(resContainerPut.Status()) {
		return resultStatusErrorResponse()
	}

	containerID := *resContainerPut.ID()
	return response(reflect.TypeOf(containerID), containerID.String())
}

//export GetContainer
func GetContainer(clientID *C.char, containerID *C.char) C.pointerResponse {
	ctx := context.Background()

	id, err := getContainerIDFromC(containerID)
	if err != nil {
		return pointerResponseError(err.Error())
	}

	var prmContainerGet neofsclient.PrmContainerGet
	prmContainerGet.SetContainer(*id)
	//prmContainerGet.WithXHeaders()

	neofsClient, err := getClient(clientID)
	if err != nil {
		return pointerResponseClientError()
	}
	neofsClient.mu.Lock()
	resContainerGet, err := neofsClient.client.ContainerGet(ctx, prmContainerGet)
	neofsClient.mu.Unlock()

	if err != nil {
		return pointerResponseError(err.Error())
	}
	if !apistatus.IsSuccessful(resContainerGet.Status()) {
		return resultStatusErrorResponsePointer()
	}

	cnr := resContainerGet.Container()
	var v2 v2container.Container
	cnr.WriteToV2(&v2)
	if err != nil {
		return pointerResponseError(err.Error())
	}
	bytes := v2.StableMarshal(nil)
	return pointerResponse(reflect.TypeOf(v2), bytes)
}

//export DeleteContainer
func DeleteContainer(clientID *C.char, containerID *C.char) C.pointerResponse {
	id, err := getContainerIDFromC(containerID)
	if err != nil {
		return pointerResponseError(err.Error())
	}

	neofsClient, err := getClient(clientID)
	if err != nil {
		return pointerResponseClientError()
	}
	neofsClient.mu.Lock()
	return deleteContainer(neofsClient, id, nil)
}

////export DeleteContainerWithinSession
//func DeleteContainerWithinSession(clientID *C.char, containerID *C.char, sessionToken *C.char) C.pointerResponse {
//	id, err := getContainerIDFromC(containerID)
//	if err != nil {
//		return pointerResponseError(err.Error())
//	}
//
//	tok, err := getSessionTokenFromC(sessionToken)
//	if err != nil {
//		return pointerResponseError(err.Error())
//	}
//
//	neofsClient, err := getClient(clientID)
//	if err != nil {
//		return pointerResponseClientError()
//	}
//	neofsClient.mu.Lock()
//	return deleteContainer(neofsClient, id, tok)
//}

func deleteContainer(neofsClient *NeoFSClient, containerID *cid.ID, sessionToken *session.Container) C.pointerResponse {
	ctx := context.Background()

	var prmContainerDelete neofsclient.PrmContainerDelete
	prmContainerDelete.SetContainer(*containerID)
	if sessionToken != nil {
		prmContainerDelete.WithinSession(*sessionToken)
	}
	//prmContainerDelete.WithXHeaders()

	resContainerDelete, err := neofsClient.client.ContainerDelete(ctx, prmContainerDelete)
	neofsClient.mu.Unlock()
	if err != nil {
		pointerResponseError(err.Error())
	}

	if !apistatus.IsSuccessful(resContainerDelete.Status()) {
		return resultStatusErrorResponsePointer()
	}

	return pointerResponseBoolean(true)
}

////export ListContainer
//func ListContainer(clientID *C.char, ownerPubKey *C.char) *C.response {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return responseClientError()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	var prmContainerList neofsclient.PrmContainerList
//	prmContainerList.SetAccount(getOwnerID(ownerPubKey))
//
//	resContainerList, err := cli.client.ContainerList(ctx, prmContainerList)
//	cli.mu.RUnlock()
//	if err != nil {
//		return responseError("could not get container list")
//	}
//	if !apistatus.IsSuccessful(resContainerList.Status()) {
//		return resultStatusErrorResponse()
//	}
//	containerIDs := resContainerList.Containers()
//	return response("ContainerList", containerIDs[0]) // how return []cid.ID
//}

////export SetExtendedACL
//func SetExtendedACL(clientID *C.char, v2Table *C.char) C.pointerResponse {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return pointerResponseClientError()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	table, err := getTableFromV2(v2Table)
//	if err != nil {
//		return pointerResponseError(err.Error())
//	}
//	var prmContainerSetEACL neofsclient.PrmContainerSetEACL
//	prmContainerSetEACL.SetTable(*table)
//
//	resContainerSetEACL, err := cli.client.ContainerSetEACL(ctx, prmContainerSetEACL)
//	cli.mu.RUnlock()
//	if err != nil {
//		return pointerResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(resContainerSetEACL.Status()) {
//		return resultStatusErrorResponsePointer()
//	}
//	boolean := []byte{1}
//	return pointerResponse(reflect.TypeOf(boolean), boolean)
//}

////export GetExtendedACL
//func GetExtendedACL(clientID *C.char, v2ContainerID *C.char) C.pointerResponse {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return pointerResponseClientError()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	containerID, err := getV2ContainerIDFromC(v2ContainerID)
//	if err != nil {
//		return pointerResponseError(err.Error())
//	}
//	var prmContainerEACL neofsclient.PrmContainerEACL
//	prmContainerEACL.SetContainer(*containerID)
//
//	cnrResponse, err := cli.client.ContainerEACL(ctx, prmContainerEACL)
//	cli.mu.RUnlock()
//	if err != nil {
//		return pointerResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(cnrResponse.Status()) {
//		return resultStatusErrorResponsePointer()
//	}
//	table := cnrResponse.Table()
//	tableBytes, err := cnrResponse.Table().Marshal()
//	if err != nil {
//		return pointerResponseError("could not marshal eacl table")
//	}
//	return pointerResponse(reflect.TypeOf(table), tableBytes)
//}

////export AnnounceUsedSpace
//func AnnounceUsedSpace(clientID *C.char, announcements *C.char) C.pointerResponse {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return pointerResponseClientError()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	ann := getAnnouncementsFromV2(announcements)
//
//	var prmContainerAnnounceSpace neofsclient.PrmAnnounceSpace
//	prmContainerAnnounceSpace.SetValues(ann)
//
//	resContainerAnnounceUsedSpace, err := cli.client.ContainerAnnounceUsedSpace(ctx, prmContainerAnnounceSpace)
//	cli.mu.RUnlock()
//	if err != nil {
//		return pointerResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(resContainerAnnounceUsedSpace.Status()) {
//		return resultStatusErrorResponsePointer()
//	}
//	boolean := []byte{1}
//	return pointerResponse(reflect.TypeOf(boolean), boolean)
//}

//endregion container
//region helper

func getContainerFromC(v2Container *C.char) (*container.Container, error) {
	v2cnr := new(v2container.Container)
	err := v2cnr.UnmarshalJSON([]byte(C.GoString(v2Container)))
	if err != nil {
		return nil, err
	}
	//v2cnr.SetHomomorphicHashingState()

	var cnr container.Container
	err = cnr.ReadFromV2(*v2cnr)
	if err != nil {
		return nil, err
	}
	return &cnr, nil
}

func getContainerIDFromC(containerID *C.char) (*cid.ID, error) {
	id := new(cid.ID)
	err := id.DecodeString(C.GoString(containerID))
	if err != nil {
		return nil, err
	}
	return id, nil
}

//func getSessionTokenFromC(sessionToken *C.char) (*session.Container, error) {
//	token := new(session.Container)
//
//	err := token.Unmarshal([]byte(C.GoString(sessionToken)))
//	if err != nil {
//		return nil, err
//	}
//	return token, nil
//}

//endregion helper
//region container old

////export PutContainerBasic
//func PutContainerBasic(key *C.char) *C.char {
//	TESTNET := "grpcs://st01.testnet.fs.neo.org:8082"
//	// create client from parameter
//	//ctx := context.TODO()
//	ctx := context.Background()
//	//walletCli, err := client.New(ctx, "http://seed1t4.neo.org:2332", client.Options{}) // get Neo endpoint from parameter
//	//if err != nil {
//	//	return fmt.Errorf("can't create wallet client: %w", err)
//	//}
//
//	privateKey := keys.PrivateKey{PrivateKey: *getECDSAPrivKey(key)}
//	ownerAcc := wallet.NewAccountFromPrivateKey(&privateKey)
//	fsCli, err := neofsclient.New(
//		neofsclient.WithDefaultPrivateKey(&privateKey.PrivateKey),
//		neofsclient.WithURIAddress(TESTNET, nil), // get NeoFS endpoint from parameter
//		neofsclient.WithNeoFSErrorParsing(),
//	)
//	if err != nil {
//		panic(fmt.Errorf("can't create neofs client: %w", err))
//	}
//
//	//	create container from parameter
//	//	required:
//	//	o	create placement policy
//	//	x	access to private key
//	//	o	set permissions
//	//	o	neofs client
//
//	ownerID := getOwnerIDFromAccount(ownerAcc)
//
//	placementPolicy := netmap.NewPlacementPolicy() // get placement policy from string
//
//	permissions := acl.PublicBasicRule
//	//acl.ParseBasicACL(aclString) // get acl from string argument
//
//	cnr := container.New(
//		container.WithPolicy(placementPolicy),
//		container.WithOwnerID(ownerID),
//		container.WithCustomBasicACL(permissions),
//	)
//
//	//attributes := container.Attributes{} // get attributes from string argument
//	//cnr.SetAttributes(attributes)
//
//	var prmContainerPut neofsclient.PrmContainerPut
//	prmContainerPut.SetContainer(*cnr)
//
//	cnrResponse, err := fsCli.ContainerPut(ctx, prmContainerPut)
//	if err != nil {
//		panic(err)
//	}
//
//	containerID := cnrResponse.ID().String()
//	cstr := C.CString(containerID)
//	return cstr
//}

// old code
////export NewContainerPutRequest
//func NewContainerPutRequest(key *C.char, v2Container *C.char) *C.char {
//	privKey := getECDSAPrivKey(key)
//
//	cnr, err := getV2ContainerFromC(v2Container)
//	if err != nil {
//		panic("could not get container from v2")
//	}
//	if cnr.Version() == nil {
//		cnr.SetVersion(version.Current())
//	}
//	_, err = cnr.NonceUUID()
//	if err != nil {
//		rand, err := uuid.NewRandom()
//		if err != nil {
//			panic("can't create new random " + err.Error())
//		}
//		cnr.SetNonceUUID(rand)
//	}
//	if cnr.BasicACL() == 0 {
//		cnr.SetBasicACL(acl.PrivateBasicRule)
//	}
//
//	// form request body
//	reqBody := new(v2container.PutRequestBody)
//	reqBody.SetContainer(cnr.ToV2())
//
//	// sign cnr
//	signWrapper := signature.StableMarshalerWrapper{SM: reqBody.GetContainer()}
//	err = sigutil.SignDataWithHandler(privKey, signWrapper, func(key []byte, sig []byte) {
//		containerSignature := new(refs.Signature)
//		containerSignature.SetKey(key)
//		containerSignature.SetSign(sig)
//		reqBody.SetSignature(containerSignature)
//	}, sigutil.SignWithRFC6979())
//	die(err)
//
//	// form meta header
//	var meta v2session.RequestMetaHeader
//	meta.SetSessionToken(cnr.SessionToken().ToV2())
//
//	// form request
//	var req v2container.PutRequest
//	req.SetBody(reqBody)
//
//	// Prepare Meta Header
//	// TODO: Check meta header params and set them accordingly
//	// 	i.e., ttl, version, network magic
//	req.SetMetaHeader(&meta)
//
//	prepareMetaHeader(&req)
//
//	pr := getRequestToSigned(&req)
//
//	err = signature.SignServiceMessage(privKey, pr)
//	die(err)
//
//	jsonAfter, err := message.MarshalJSON(pr)
//	die(err)
//
//	cstr := C.CString(string(jsonAfter))
//
//	return cstr
//}

//endregion container old
