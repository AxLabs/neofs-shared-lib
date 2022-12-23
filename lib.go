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
	"bytes"
	"fmt"
	"github.com/AxLabs/neofs-api-shared-lib/accounting"
	"github.com/AxLabs/neofs-api-shared-lib/client"
	"github.com/AxLabs/neofs-api-shared-lib/container"
	"github.com/AxLabs/neofs-api-shared-lib/netmap"
	"github.com/AxLabs/neofs-api-shared-lib/object"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	"unsafe"
)

// region client

//export CreateClient
func CreateClient(privateKey *C.char, neofsEndpoint *C.char) C.pointerResponse {
	privKey := GetECDSAPrivKey(privateKey)
	endpoint := toGoString(neofsEndpoint)
	return responseToC(client.CreateClient(privKey, endpoint))
}

// endregion client
// region accounting

//export GetBalance
func GetBalance(clientID *C.char, publicKey *C.char) C.pointerResponse {
	c, err := GetClient(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	key, err := UserIDFromPublicKey(publicKey)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(accounting.GetBalance(c, key))
}

// endregion accounting
// region netmap

//export GetEndpoint
func GetEndpoint(clientID *C.char) C.pointerResponse {
	c, err := GetClient(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(netmap.GetEndpoint(c))
}

//export GetNetworkInfo
func GetNetworkInfo(clientID *C.char) C.pointerResponse {
	c, err := GetClient(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(netmap.GetNetworkInfo(c))
}

// endregion netmap
// region container

//export PutContainer
func PutContainer(clientID *C.char, v2Container *C.char) C.response {
	sdkContainer, err := getContainerFromC(v2Container)
	if err != nil {
		return stringResponseToC(response.StringError(err))
	}

	c, err := GetClient(clientID)
	if err != nil {
		return stringResponseToC(response.StringError(err))
	}
	return stringResponseToC(container.PutContainer(c, sdkContainer))
}

//export GetContainer
func GetContainer(clientID *C.char, containerID *C.char) C.pointerResponse {
	cid, err := getContainerIDFromC(containerID)
	if err != nil {
		return responseToC(response.Error(err))
	}

	c, err := GetClient(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(container.GetContainer(c, cid))
}

//export DeleteContainer
func DeleteContainer(clientID *C.char, containerID *C.char) C.pointerResponse {
	cid, err := getContainerIDFromC(containerID)
	if err != nil {
		return responseToC(response.Error(err))
	}

	c, err := GetClient(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	container.DeleteContainer(c, cid)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(response.NewBoolean(true))
}

//export ListContainer
func ListContainer(clientID *C.char, ownerPubKey *C.char) C.pointerResponse {
	userID, err := UserIDFromPublicKey(ownerPubKey)
	if err != nil {
		return responseToC(response.Error(err))
	}
	c, err := GetClient(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(container.ListContainer(c, userID))
}

// endregion container
// region object

////export UploadFile
//func UploadFile(clientID *C.char, fileBytes *C.char) C.pointerResponse {
//	c, err := GetClient(clientID)
//	if err != nil {
//		return responseToC(response.Error(err))
//	}
//	prm := UploadFilePrm{}
//	err = json.Unmarshal(fileBytes, &prm)
//	if err != nil {
//		return responseToC(response.Error(err))
//	}
//	return responseToC(object.UploadFile(c, prm))
//}
//
//type UploadFilePrm struct {
//	prm1 string
//	prm2 string
//}

//export CreateObjectWithoutAttributes
func CreateObjectWithoutAttributes(clientID *C.char, containerID *C.char, fileBytes unsafe.Pointer, fileSize C.int,
	sessionSignerPrivKey *C.char) C.response {
	return CreateObject(clientID, containerID, fileBytes, fileSize, sessionSignerPrivKey, nil, nil)
}

// seems to work
//export CreateObject
func CreateObject(clientID *C.char, containerID *C.char, fileBytes unsafe.Pointer, fileSize C.int, sessionSignerPrivKey *C.char,
	attributeKey *C.char, attributeValue *C.char) C.response {

	readBytes := C.GoBytes(fileBytes, fileSize)
	fmt.Println(string(readBytes))

	c, err := GetClient(clientID)
	if err != nil {
		return stringResponseToC(response.StringError(err))
	}
	reader := bytes.NewReader(readBytes)
	fmt.Println("reader initialized")
	privKey := GetECDSAPrivKey(sessionSignerPrivKey)
	cid, err := getContainerIDFromC(containerID)
	if err != nil {
		return stringResponseToC(response.StringError(err))
	}
	var attributes [][2]string
	if attributeKey != nil {
		key := C.GoString(attributeKey)
		value := C.GoString(attributeValue)
		attributes = append(attributes, [2]string{key, value})
		fmt.Println("attributes initialized")
	}

	return stringResponseToC(object.CreateObject(c, *cid, *privKey, attributes, reader))
}

//ReadObject(neofsClient *client.NeoFSClient, containerID cid.ID, objectID oid.ID,
//signer ecdsa.PrivateKey) *response.PointerResponse {

// object not found
//export ReadObject
func ReadObject(clientID *C.char, containerID *C.char, objectID *C.char, signer *C.char) C.pointerResponse {
	c, err := GetClient(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	cid, err := getContainerIDFromC(containerID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	oid, err := getObjectIDFromC(objectID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	privKey := GetECDSAPrivKey(signer)

	return responseToC(object.ReadObject(c, *cid, *oid, *privKey))
}

// endregion object
