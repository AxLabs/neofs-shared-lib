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
)

/*
----Accounting----
Balance
*/

//export GetBalance
func GetBalance(clientID *C.char, publicKey *C.char) C.response {
	cli, err := getClient(clientID)
	if err != nil {
		return cResponseError(err.Error())
	}
	cli.mu.RLock()
	var prmBalanceGet neofsCli.PrmBalanceGet
	ownerID := getOwnerID(getPubKey(publicKey))
	prmBalanceGet.SetAccount(ownerID)
	ctx := context.Background()
	resBalanceGet, err := cli.client.BalanceGet(ctx, prmBalanceGet)
	cli.mu.RUnlock()

	if err != nil {
		return cResponseError("could not get endpoint info")
	}
	status := resBalanceGet.Status()
	if !apistatus.IsSuccessful(status) {
		return cResponseErrorStatus()
	}
	amount := resBalanceGet.Amount()
	if amount == nil {
		return cResponseError("could not get balance")
	}
	json, err := amount.MarshalJSON()
	if err != nil {
		return cResponseError("could not marshal network info of endpoint")
	}
	return C.response{C.CString("GetBalance"), (*C.char)(C.CBytes(json))}
}
