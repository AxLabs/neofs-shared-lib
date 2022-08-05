package main

/*
#include <stdlib.h>

#ifndef RESPONSE_H
#define RESPONSE_H
#include "response.h"
#endif
*/
import "C"

/*
----Session----
Create
*/
//export CreateSession
func CreateSession(clientID *C.char, sessionExpiration *C.ulonglong) {}

//func CreateSession(clientID *C.char, sessionExpiration *C.ulonglong) C.response {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return clientErrorResponse()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	exp := sessionExpiration // uint
//	var prmSessionCreate neofsCli.PrmSessionCreate
//	prmSessionCreate.SetExp(exp)
//
//	resSessionCreate, err := cli.client.SessionCreate(ctx, prmSessionCreate)
//	cli.mu.RUnlock()
//	if err != nil {
//		return errorResponse(err.Error())
//	}
//	if !apistatus.IsSuccessful(resSessionCreate.Status()) {
//		return resultStatusErrorResponse()
//	}
//	sessionID := resSessionCreate.ID()
//	sessionPublicKey := resSessionCreate.PublicKey()
//	return newResponse("CreateSession", ) // handle method with two return values
//}
