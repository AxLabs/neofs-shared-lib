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
	neofsCli "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/object"
	oid "github.com/nspcc-dev/neofs-sdk-go/object/id"
	"reflect"
)

/*
----Object----
Get
Put
Delete
Head
Search
GetRange
GetRangeHash
*/

////export GetObjectInit
//func GetObjectInit(containerID *C.char, neofsEndpoint *C.char, key *C.char) *C.char {
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	// Parse the container
//	id := getContainerIDFromV2(containerID)
//
//	var prmObjectGet neofsCli.PrmObjectGet
//	prmObjectGet.FromContainer()
//
//	response, err := fsCli.ObjectGetInit(ctx, prmObjectGet)
//	if err != nil {
//		panic(err)
//	}
//
//	response.Read()
//	return C.CString(containerJson)
//}

////export ReadObject
//func ReadObject(objectReader *C.char) {
//
//}

////export PutObject
//func PutObject(neofsEndpoint *C.char, key *C.char) *C.char {
//	fsCli, err := getClient(key)
//	ctx := context.Background()
//
//	var prmObjectPutInit neofsCli.PrmObjectPutInit
//	response, err := fsCli.client.ObjectPutInit(ctx, prmObjectPutInit)
//	if err != nil {
//		panic(err)
//	}
//
//	//response.WritePayloadChunk()
//	//response.
//	return C.CString(containerJson)
//}

//export DeleteObject
func DeleteObject(clientID *C.char, v2ContainerID *C.char, v2ObjectID *C.char, v2SessionToken *C.char, v2BearerToken *C.char) C.response {
	cli, err := getClient(clientID)
	if err != nil {
		return clientErrorResponse()
	}
	cli.mu.RLock()
	ctx := context.Background()
	containerID, err := getContainerIDFromV2(v2ContainerID)
	if err != nil {
		return errorResponse(err.Error())
	}
	objectID, err := getObjectIDFromV2(v2ObjectID)
	if err != nil {
		return errorResponse(err.Error())
	}
	sessionToken, err := getSessionTokenFromV2(v2SessionToken)
	if err != nil {
		return errorResponse(err.Error())
	}
	bearerToken, err := getBearerTokenFromV2(v2BearerToken)
	if err != nil {
		return errorResponse(err.Error())
	}
	var prmObjectDelete neofsCli.PrmObjectDelete
	prmObjectDelete.FromContainer(*containerID)
	prmObjectDelete.ByID(*objectID)
	prmObjectDelete.WithinSession(*sessionToken)
	prmObjectDelete.WithBearerToken(*bearerToken)

	resObjectDelete, err := cli.client.ObjectDelete(ctx, prmObjectDelete)
	cli.mu.RUnlock()
	if err != nil {
		return errorResponse(err.Error())
	}
	if !apistatus.IsSuccessful(resObjectDelete.Status()) {
		return resultStatusErrorResponse()
	}
	readTombStoneID := new(oid.ID)
	tombstoneRead := resObjectDelete.ReadTombstoneID(readTombStoneID)
	if !tombstoneRead {
		return errorResponse("could not read object's tombstone")
	}
	json, err := readTombStoneID.MarshalJSON()
	if err != nil {
		return errorResponse(err.Error())
	}
	return newResponse(reflect.TypeOf(tombstoneRead), json)
}

//export GetObjectHead
func GetObjectHead(clientID *C.char, v2ContainerID *C.char, v2ObjectID *C.char, v2SessionToken *C.char, v2BearerToken *C.char) C.response {
	cli, err := getClient(clientID)
	if err != nil {
		return clientErrorResponse()
	}
	cli.mu.RLock()
	ctx := context.Background()
	containerID, err := getContainerIDFromV2(v2ContainerID)
	if err != nil {
		return errorResponse(err.Error())
	}
	objectID, err := getObjectIDFromV2(v2ObjectID)
	if err != nil {
		return errorResponse(err.Error())
	}
	sessionToken, err := getSessionTokenFromV2(v2SessionToken)
	if err != nil {
		return errorResponse(err.Error())
	}
	bearerToken, err := getBearerTokenFromV2(v2BearerToken)
	if err != nil {
		return errorResponse(err.Error())
	}
	var prmObjectHead neofsCli.PrmObjectHead
	prmObjectHead.FromContainer(*containerID)
	prmObjectHead.ByID(*objectID)
	prmObjectHead.WithinSession(*sessionToken)
	prmObjectHead.WithBearerToken(*bearerToken)

	resObjectHead, err := cli.client.ObjectHead(ctx, prmObjectHead)
	cli.mu.RUnlock()
	if err != nil {
		panic(err)
	}
	if !apistatus.IsSuccessful(resObjectHead.Status()) {
		return resultStatusErrorResponse()
	}
	dst := new(object.Object)
	resObjectHead.ReadHeader(dst)
	json, err := dst.MarshalJSON()
	if err != nil {
		return errorResponse(err.Error())
	}
	return newResponse(reflect.TypeOf(dst), json)
}

//export SearchObject
func SearchObject(clientID *C.char, v2ContainerID *C.char, v2SessionToken *C.char, v2BearerToken *C.char, v2Filters *C.char) {
}

//func SearchObject(clientID *C.char, v2ContainerID *C.char, v2SessionToken *C.char, v2BearerToken *C.char, v2Filters *C.char) C.response {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return clientErrorResponse()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//
//	containerID, err := getContainerIDFromV2(v2ContainerID)
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//	sessionToken, err := getSessionTokenFromV2(v2SessionToken)
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//	bearerToken, err := getBearerTokenFromV2(v2BearerToken)
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//	filters, err := getFiltersFromV2(v2Filters)
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//	var prmObjectSearch neofsCli.PrmObjectSearch
//	prmObjectSearch.InContainer(*containerID)
//	prmObjectSearch.WithinSession(*sessionToken)
//	prmObjectSearch.WithBearerToken(*bearerToken)
//	prmObjectSearch.SetFilters(*filters)
//	//prmObjectSearch.MarkLocal()
//
//	resObjectSearchInit, err := cli.client.ObjectSearchInit(ctx, prmObjectSearch)
//	cli.mu.RUnlock()
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//
//	//resObjectSearchInit.UseKey()
//	//resObjectSearchInit.Read()
//	//resObjectSearchInit.Close()
//
//	read, b := resObjectSearchInit.Read()
//	return newResponse("SearchObject", read)
//}

//export GetRange
func GetRange(clientID *C.char, v2ContainerID *C.char, v2ObjectID *C.char, v2SessionToken *C.char, v2BearerToken *C.char, length *C.char,
	offset *C.char) {
}

//func GetRange(clientID *C.char, v2ContainerID *C.char, v2ObjectID *C.char, v2SessionToken *C.char, v2BearerToken *C.char, length *C.char,
//	offset *C.char) C.response {
//
//	cli, err := getClient(clientID)
//	if err != nil {
//		return clientErrorResponse()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//
//	containerID, err := getContainerIDFromV2(v2ContainerID)
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//	objectID, err := getObjectIDFromV2(v2ObjectID)
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//	sessionToken, err := getSessionTokenFromV2(v2SessionToken)
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//	bearerToken, err := getBearerTokenFromV2(v2BearerToken)
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//
//	var prmObjectRange neofsCli.PrmObjectRange
//	prmObjectRange.FromContainer(*containerID)
//	prmObjectRange.ByID(*objectID)
//	prmObjectRange.WithinSession(*sessionToken)
//	prmObjectRange.WithBearerToken(*bearerToken)
//	prmObjectRange.SetLength(length)
//	prmObjectRange.SetOffset(offset)
//
//	response, err := cli.client.ObjectRangeInit(ctx, prmObjectRange)
//	cli.mu.RUnlock()
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//
//	response.Read()
//	return newResponse("GetRange", )
//}

//export GetRangeHash
func GetRangeHash(clientID *C.char, v2ContainerID *C.char) {
}
