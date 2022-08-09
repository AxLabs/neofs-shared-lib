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
	"reflect"
	"strconv"
)

/*
----Session----
Create
*/
//export CreateSession
//func CreateSession(clientID *C.char, sessionExpiration *C.ulonglong) {}

func CreateSession(clientID *C.char, sessionExpiration *C.char) C.pointerResponse {
	cli, err := getClient(clientID)
	if err != nil {
		return pointerResponseClientError()
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmSessionCreate neofsCli.PrmSessionCreate
	exp, err := strconv.ParseUint(C.GoString(sessionExpiration), 10, 64)
	if err != nil {
		return pointerResponseError("could not parse session expiration to uint64")
	}
	prmSessionCreate.SetExp(exp)

	resSessionCreate, err := cli.client.SessionCreate(ctx, prmSessionCreate)
	cli.mu.RUnlock()
	if err != nil {
		return pointerResponseError(err.Error())
	}
	if !apistatus.IsSuccessful(resSessionCreate.Status()) {
		return resultStatusErrorResponsePointer()
	}
	sessionID := resSessionCreate.ID()
	//sessionPublicKey := resSessionCreate.PublicKey()
	return pointerResponse(reflect.TypeOf(sessionID), sessionID) // handle method with two return values
}

//export CreateSessionPubKey
func CreateSessionPubKey(clientID *C.char, sessionExpiration *C.char) C.pointerResponse {
	cli, err := getClient(clientID)
	if err != nil {
		return pointerResponseClientError()
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmSessionCreate neofsCli.PrmSessionCreate
	exp, err := strconv.ParseUint(C.GoString(sessionExpiration), 10, 64)
	if err != nil {
		return pointerResponseError("could not parse session expiration to uint64")
	}
	prmSessionCreate.SetExp(exp)

	resSessionCreate, err := cli.client.SessionCreate(ctx, prmSessionCreate)
	cli.mu.RUnlock()
	if err != nil {
		return pointerResponseError(err.Error())
	}
	if !apistatus.IsSuccessful(resSessionCreate.Status()) {
		return resultStatusErrorResponsePointer()
	}
	sessionPublicKey := resSessionCreate.PublicKey()
	return pointerResponse(reflect.TypeOf(sessionPublicKey), sessionPublicKey) // handle method with two return values
}
