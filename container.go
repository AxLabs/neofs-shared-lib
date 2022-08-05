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
	"fmt"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	"github.com/nspcc-dev/neofs-sdk-go/acl"
	neofsCli "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	cid "github.com/nspcc-dev/neofs-sdk-go/container/id"
	"github.com/nspcc-dev/neofs-sdk-go/netmap"
	"github.com/nspcc-dev/neofs-sdk-go/version"
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
func PutContainer(clientID *C.char, v2Container *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return clientErrorResponsePointer()
	}
	cli.mu.RLock()
	ctx := context.Background()
	cnr, err := getContainerFromV2(v2Container)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	// Overwrites potentially set container version
	cnr.SetVersion(version.Current())
	// The following are expected to be set within the provided container parameter
	//  - placement policy
	//  - permissions
	//  - attributes
	var prmContainerPut neofsCli.PrmContainerPut
	prmContainerPut.SetContainer(*cnr)

	resContainerPut, err := cli.client.ContainerPut(ctx, prmContainerPut)
	if err != nil {
		return errorResponsePointer("could not put container")
	}
	if !apistatus.IsSuccessful(resContainerPut.Status()) {
		return resultStatusErrorResponsePointer()
	}
	containerID := resContainerPut.ID()
	json, err := containerID.MarshalJSON()
	if err != nil {
		return errorResponsePointer("could not marshal container put response")
	}
	return newResponsePointer(reflect.TypeOf(containerID), json)
}

//export GetContainer
func GetContainer(clientID *C.char, v2ContainerID *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return clientErrorResponsePointer()
	}
	cli.mu.RLock()
	ctx := context.Background()
	id, err := getContainerIDFromV2(v2ContainerID)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	var prmContainerGet neofsCli.PrmContainerGet
	var prmContainerPut neofsCli.PrmContainerPut
	prmContainerPut.SetContainer(container.Container{})
	prmContainerGet.SetContainer(*id)

	resContainerGet, err := cli.client.ContainerGet(ctx, prmContainerGet)
	cli.mu.RUnlock()

	if err != nil {
		return errorResponsePointer("could not get container")
	}
	if !apistatus.IsSuccessful(resContainerGet.Status()) {
		return resultStatusErrorResponsePointer()
	}
	ctr := resContainerGet.Container()
	containerJson, err := ctr.MarshalJSON()
	if err != nil {
		return errorResponsePointer("could not marshal container put response")
	}
	return newResponsePointer(reflect.TypeOf(ctr), containerJson)
}

//export DeleteContainer
func DeleteContainer(clientID *C.char, v2ContainerID *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return clientErrorResponsePointer()
	}
	cli.mu.RLock()
	ctx := context.Background()
	id, err := getContainerIDFromV2(v2ContainerID)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	var prmContainerDelete neofsCli.PrmContainerDelete
	prmContainerDelete.SetContainer(*id)
	//prmContainerGet.SetSessionToken()

	resContainerDelete, err := cli.client.ContainerDelete(ctx, prmContainerDelete)
	if err != nil {
		panic(err)
	}

	if !apistatus.IsSuccessful(resContainerDelete.Status()) {
		return resultStatusErrorResponsePointer()
	}
	boolean := []byte{1}
	return newResponsePointer(reflect.TypeOf(resContainerDelete), boolean) // handle methods without return value
}

//export ListContainer
func ListContainer(clientID *C.char, ownerPubKey *C.char) {}

//func ListContainer(clientID *C.char, ownerPubKey *C.char) *C.response {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return clientErrorResponse()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	var prmContainerList neofsCli.PrmContainerList
//	prmContainerList.SetAccount(getOwnerID(ownerPubKey))
//
//	resContainerList, err := cli.client.ContainerList(ctx, prmContainerList)
//	cli.mu.RUnlock()
//	if err != nil {
//		return errorResponse("could not get container list")
//	}
//	if !apistatus.IsSuccessful(resContainerList.Status()) {
//		return resultStatusErrorResponse()
//	}
//	containerIDs := resContainerList.Containers()
//	return newResponse("ContainerList", containerIDs[0]) // how return []cid.ID
//}

//export SetExtendedACL
func SetExtendedACL(clientID *C.char, v2Table *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return clientErrorResponsePointer()
	}
	cli.mu.RLock()
	ctx := context.Background()
	table, err := getTableFromV2(v2Table)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	var prmContainerSetEACL neofsCli.PrmContainerSetEACL
	prmContainerSetEACL.SetTable(*table)

	resContainerSetEACL, err := cli.client.ContainerSetEACL(ctx, prmContainerSetEACL)
	cli.mu.RUnlock()
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	if !apistatus.IsSuccessful(resContainerSetEACL.Status()) {
		return resultStatusErrorResponsePointer()
	}
	boolean := []byte{1}
	return newResponsePointer(reflect.TypeOf(boolean), boolean)
}

//export GetExtendedACL
func GetExtendedACL(clientID *C.char, v2ContainerID *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return clientErrorResponsePointer()
	}
	cli.mu.RLock()
	ctx := context.Background()
	containerID, err := getContainerIDFromV2(v2ContainerID)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	var prmContainerEACL neofsCli.PrmContainerEACL
	prmContainerEACL.SetContainer(*containerID)

	cnrResponse, err := cli.client.ContainerEACL(ctx, prmContainerEACL)
	cli.mu.RUnlock()
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	if !apistatus.IsSuccessful(cnrResponse.Status()) {
		return resultStatusErrorResponsePointer()
	}
	table := cnrResponse.Table()
	containerJson, err := table.MarshalJSON()
	if err != nil {
		return errorResponsePointer("could not marshal container put response")
	}
	return newResponsePointer(reflect.TypeOf(table), containerJson)
}

//export AnnounceUsedSpace
func AnnounceUsedSpace(clientID *C.char, announcements *C.char) {}

//func AnnounceUsedSpace(clientID *C.char, announcements *C.char) C.responsePointer {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return clientErrorResponsePointer()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	ann := getAnnouncementsFromV2(announcements)
//
//	var prmContainerAnnounceSpace neofsCli.PrmAnnounceSpace
//	prmContainerAnnounceSpace.SetValues(ann)
//
//	resContainerAnnounceUsedSpace, err := cli.client.ContainerAnnounceUsedSpace(ctx, prmContainerAnnounceSpace)
//	cli.mu.RUnlock()
//	if err != nil {
//		return errorResponsePointer(err.Error())
//	}
//	if !apistatus.IsSuccessful(resContainerAnnounceUsedSpace.Status()) {
//		return resultStatusErrorResponsePointer()
//	}
//	boolean := []byte{1}
//	return newResponsePointer(reflect.TypeOf(boolean), boolean)
//}

//endregion container
//region helper

func getContainerFromV2(v2Container *C.char) (*container.Container, error) {
	sdkContainer := new(container.Container)
	str := C.GoString(v2Container)
	err := sdkContainer.UnmarshalJSON([]byte(str))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal container")
	}
	return sdkContainer, nil
}

func getContainerIDFromV2(containerID *C.char) (*cid.ID, error) {
	id := new(cid.ID)
	err := id.UnmarshalJSON([]byte(C.GoString(containerID)))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal container id")
	}
	return id, nil
}

//endregion helper
//region container old

//export PutContainerBasic
func PutContainerBasic(key *C.char) *C.char {
	TESTNET := "grpcs://st01.testnet.fs.neo.org:8082"
	// create client from parameter
	//ctx := context.TODO()
	ctx := context.Background()
	//walletCli, err := client.New(ctx, "http://seed1t4.neo.org:2332", client.Options{}) // get Neo endpoint from parameter
	//if err != nil {
	//	return fmt.Errorf("can't create wallet client: %w", err)
	//}

	privateKey := keys.PrivateKey{PrivateKey: *getECDSAPrivKey(key)}
	ownerAcc := wallet.NewAccountFromPrivateKey(&privateKey)
	fsCli, err := neofsCli.New(
		neofsCli.WithDefaultPrivateKey(&privateKey.PrivateKey),
		neofsCli.WithURIAddress(TESTNET, nil), // get NeoFS endpoint from parameter
		neofsCli.WithNeoFSErrorParsing(),
	)
	if err != nil {
		panic(fmt.Errorf("can't create neofs client: %w", err))
	}

	//	create container from parameter
	//	required:
	//	o	create placement policy
	//	x	access to private key
	//	o	set permissions
	//	o	neofs client

	ownerID := getOwnerIDFromAccount(ownerAcc)

	placementPolicy := netmap.NewPlacementPolicy() // get placement policy from string

	permissions := acl.PublicBasicRule
	//acl.ParseBasicACL(aclString) // get acl from string argument

	cnr := container.New(
		container.WithPolicy(placementPolicy),
		container.WithOwnerID(ownerID),
		container.WithCustomBasicACL(permissions),
	)

	//attributes := container.Attributes{} // get attributes from string argument
	//cnr.SetAttributes(attributes)

	var prmContainerPut neofsCli.PrmContainerPut
	prmContainerPut.SetContainer(*cnr)

	cnrResponse, err := fsCli.ContainerPut(ctx, prmContainerPut)
	if err != nil {
		panic(err)
	}

	containerID := cnrResponse.ID().String()
	cstr := C.CString(containerID)
	return cstr
}

// old code
////export NewContainerPutRequest
//func NewContainerPutRequest(key *C.char, v2Container *C.char) *C.char {
//	privKey := getECDSAPrivKey(key)
//
//	cnr, err := getContainerFromV2(v2Container)
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
