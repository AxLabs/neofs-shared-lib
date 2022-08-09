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
)

/*
----Accounting----
Balance
*/

//export GetBalance
func GetBalance(clientID *C.char, publicKey *C.char) C.pointerResponse {
	cli, err := getClient(clientID)
	if err != nil {
		return pointerResponseError(err.Error())
	}
	cli.mu.RLock()
	var prmBalanceGet neofsCli.PrmBalanceGet
	ownerID := getOwnerID(getPubKey(publicKey))
	prmBalanceGet.SetAccount(ownerID)
	ctx := context.Background()
	resBalanceGet, err := cli.client.BalanceGet(ctx, prmBalanceGet)
	cli.mu.RUnlock()

	if err != nil {
		return pointerResponseError("could not get endpoint info")
	}
	status := resBalanceGet.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	amount := resBalanceGet.Amount()
	if amount == nil {
		return pointerResponseError("could not get balance")
	}
	bytes, err := amount.Marshal()
	if err != nil {
		return pointerResponseError("could not marshal balance amount")
	}
	return pointerResponse(reflect.TypeOf(amount), bytes)
}
