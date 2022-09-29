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
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/encoding/address"
	v2container "github.com/nspcc-dev/neofs-api-go/v2/container"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	cid "github.com/nspcc-dev/neofs-sdk-go/container/id"
	"github.com/nspcc-dev/neofs-sdk-go/user"
	"math/big"
)

func toGoString(cString *C.char) string {
	return C.GoString(cString)
}

func responseToC(p *response.PointerResponse) C.pointerResponse {
	return C.pointerResponse{
		C.CString(p.GetTypeString()),
		C.int(len(p.GetData())),
		(*C.char)(C.CBytes(p.GetData())),
	}
}

func stringResponseToC(p *response.StringResponse) C.response {
	return C.response{
		C.CString(p.GetTypeString()),
		C.CString(p.GetValue()),
	}
}

func GetECDSAPrivKey(key *C.char) *ecdsa.PrivateKey {
	keyStr := C.GoString(key)
	bytes, _ := hex.DecodeString(keyStr)
	k := new(big.Int)
	k.SetBytes(bytes)
	priv := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	priv.PublicKey.Curve = curve
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())
	return priv
}

func UserIDFromPublicKey(publicKey *C.char) (*user.ID, error) {
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

func uuidToGo(clientID *C.char) (*uuid.UUID, error) {
	id, err := uuid.Parse(C.GoString(clientID))
	if err != nil {
		return nil, fmt.Errorf("could not parse provided client id. id was " + C.GoString(clientID))
	}
	return &id, nil
}

func getContainerFromC(v2Container *C.char) (*container.Container, error) {
	v2cnr := new(v2container.Container)
	err := v2cnr.UnmarshalJSON([]byte(C.GoString(v2Container)))
	if err != nil {
		return nil, err
	}
	//v2cnr.SetHomomorphicHashingState()

	var cnr container.Container
	err = cnr.ReadFromV2(*v2cnr)
	if err != nil {
		return nil, err
	}
	return &cnr, nil
}

func getContainerIDFromC(containerID *C.char) (*cid.ID, error) {
	id := new(cid.ID)
	err := id.DecodeString(C.GoString(containerID))
	if err != nil {
		return nil, err
	}
	return id, nil
}

//func getSessionTokenFromC(sessionToken *C.char) (*session.Container, error) {
//	token := new(session.Container)
//
//	err := token.Unmarshal([]byte(C.GoString(sessionToken)))
//	if err != nil {
//		return nil, err
//	}
//	return token, nil
//}
