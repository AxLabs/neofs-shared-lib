package main

/*
#include <stdlib.h>

#ifndef RESPONSE_H
#define RESPONSE_H
#include "response.h"
#endif
*/
import "C"
import "github.com/AxLabs/neofs-api-shared-lib/client"

//func GetECDSAPrivKey(key *C.char) *ecdsa.PrivateKey {
//	keyStr := C.GoString(key)
//	bytes, _ := hex.DecodeString(keyStr)
//	k := new(big.Int)
//	k.SetBytes(bytes)
//	priv := new(ecdsa.PrivateKey)
//	curve := elliptic.P256()
//	priv.PublicKey.Curve = curve
//	priv.D = k
//	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())
//	return priv
//}

//func UserIDFromPublicKey(publicKey *C.char) (*user.ID, error) {
//	pubKey, err := keys.NewPublicKeyFromString(C.GoString(publicKey))
//	if err != nil {
//		return nil, err
//	}
//	var id user.ID
//	uint160, err := address.StringToUint160(pubKey.Address())
//	if err != nil {
//		return nil, err
//	}
//	id.SetScriptHash(uint160)
//	return &id, nil
//}

func GetClient(clientID *C.char) (*client.NeoFSClient, error) {
	id, err := uuidToGo(clientID)
	if err != nil {
		return nil, err
	}
	c, err := client.GetClient(id)
	if err != nil {
		return nil, err
	}
	return c, nil
}
