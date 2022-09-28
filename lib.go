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

//export CreateClient
func CreateClient(privateKey *C.char, neofsEndpoint *C.char) C.pointerResponse {
	privKey := GetECDSAPrivKey(privateKey)
	endpoint := toGoString(neofsEndpoint)
	return responseToC(client.CreateClient(privKey, endpoint))
}

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

//export PutContainer
func PutContainer(clientID *C.char, v2Container *C.char) C.response {
	c, err := GetClient(clientID)
	if err != nil {
		return stringResponseToC(response.StringError(err))
	}
	sdkContainer, err := getContainerFromC(v2Container)
	if err != nil {
		return stringResponseToC(response.StringError(err))
	}
	return stringResponseToC(container.PutContainer(c, sdkContainer))
}

//export GetContainer
func GetContainer(clientID *C.char, containerID *C.char) C.pointerResponse {
	c, err := GetClient(clientID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	cid, err := getContainerIDFromC(containerID)
	if err != nil {
		return responseToC(response.Error(err))
	}
	if err != nil {
		return responseToC(response.Error(err))
	}
	return responseToC(container.GetContainer(c, cid))
}
