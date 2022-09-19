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
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/encoding/address"
	v2accounting "github.com/nspcc-dev/neofs-api-go/v2/accounting"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/user"
	"reflect"
)

/*
----Accounting----
Balance
*/

//export GetBalance
func GetBalance(clientID *C.char, publicKey *C.char) C.pointerResponse {
	ctx := context.Background()

	var prmBalanceGet neofsclient.PrmBalanceGet
	id, err := userIDFromPublicKey(publicKey)
	prmBalanceGet.SetAccount(*id)

	neofsClient, err := getClient(clientID)
	if err != nil {
		return pointerResponseClientError()
	}
	neofsClient.mu.Lock()
	resBalanceGet, err := neofsClient.client.BalanceGet(ctx, prmBalanceGet)
	neofsClient.mu.Unlock()
	if err != nil {
		return pointerResponseError(err.Error())
	}

	resStatus := resBalanceGet.Status()
	if !apistatus.IsSuccessful(resStatus) {
		return resultStatusErrorResponsePointer()
	}

	amount := resBalanceGet.Amount()
	if amount == nil {
		return pointerResponseError(err.Error())
	}

	var v2 v2accounting.Decimal
	amount.WriteToV2(&v2)
	bytes := v2.StableMarshal(nil)
	return pointerResponse(reflect.TypeOf(v2), bytes)
}

func userIDFromPublicKey(publicKey *C.char) (*user.ID, error) {
	pubKey, err := keys.NewPublicKeyFromString(C.GoString(publicKey))
	if err != nil {
		return nil, err
	}
	var id user.ID
	uint160, err := address.StringToUint160(pubKey.Address())
	if err != nil {
		return nil, err
	}
	id.SetScriptHash(uint160)
	return &id, nil
}
