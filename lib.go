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
	"github.com/AxLabs/neofs-api-shared-lib/accounting"
	"github.com/AxLabs/neofs-api-shared-lib/client"
	"github.com/AxLabs/neofs-api-shared-lib/container"
	"github.com/AxLabs/neofs-api-shared-lib/netmap"
	"github.com/AxLabs/neofs-api-shared-lib/response"
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
	id, err := uuidToGo(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	key, err := UserIDFromPublicKey(publicKey)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(accounting.GetBalance(id, key))
}

// endregion accounting
// region netmap

//export GetEndpoint
func GetEndpoint(clientID *C.char) C.pointerResponse {
	id, err := uuidToGo(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(netmap.GetEndpoint(id))
}

//export GetNetworkInfo
func GetNetworkInfo(clientID *C.char) C.pointerResponse {
	id, err := uuidToGo(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(netmap.GetNetworkInfo(id))
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

// endregion object
