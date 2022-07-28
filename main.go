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
	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	v2accounting "github.com/nspcc-dev/neofs-api-go/v2/accounting"
	v2container "github.com/nspcc-dev/neofs-api-go/v2/container"
	"github.com/nspcc-dev/neofs-api-go/v2/rpc/message"
	v2session "github.com/nspcc-dev/neofs-api-go/v2/session"
	"github.com/nspcc-dev/neofs-api-go/v2/signature"
	neofsCli "github.com/nspcc-dev/neofs-sdk-go/client"
	"github.com/nspcc-dev/neofs-sdk-go/eacl"
	"github.com/nspcc-dev/neofs-sdk-go/object"
	oid "github.com/nspcc-dev/neofs-sdk-go/object/id"
	"github.com/nspcc-dev/neofs-sdk-go/owner"
	"github.com/nspcc-dev/neofs-sdk-go/reputation"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	"github.com/nspcc-dev/neofs-sdk-go/token"
	"github.com/nspcc-dev/neofs-sdk-go/version"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"math/big"
	"sync"
)

func main() {

	keyStr := "84180ac9d6eb6fba207ea4ef9d2200102d1ebeb4b9c07e2c6a738a42742e27a5"

	bytes, err := hex.DecodeString(keyStr)

	k := new(big.Int)
	k.SetBytes(bytes)

	priv := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	priv.PublicKey.Curve = curve
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())

	jsonFromJava := "{\n  \"body\": {\n    \"container\": {\n      \"version\": {\n        \"major\": 2,\n        \"minor\": 11\n      },\n      \"ownerID\": {\n        \"value\": \"A9+X+2Xt74Dy/JmsSuTv0aMMUZwH+LNiF4J4fyiBqbe1\"\n      },\n      \"nonce\": \"L4Mz1w==\",\n      \"basicACL\": 532660223,\n      \"attributes\": [{\n        \"key\": \"key\",\n        \"value\": \"val\"\n      }],\n      \"placementPolicy\": {\n        \"replicas\": [{\n          \"count\": 2\n        }],\n        \"containerBackupFactor\": 1\n      }\n    }\n  },\n  \"metaHeader\": {\n    \"version\": {\n      \"major\": 2,\n      \"minor\": 11\n    },\n    \"epoch\": \"10\",\n    \"ttl\": 1000\n  }\n}"

	pr := &v2container.PutRequest{}

	m := pr.ToGRPCMessage().(proto.Message)
	err = protojson.Unmarshal([]byte(jsonFromJava), m)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	_ = pr.FromGRPCMessage(m)

	err = signature.SignServiceMessage(priv, pr)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	jsonAfter, err := message.MarshalJSON(pr)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Println(string(jsonAfter))

}

func getOwnerIDFromAccount(acc *wallet.Account) *owner.ID {
	return owner.NewIDFromN3Account(acc)
}

func getOwnerIDFromPublicKey(pubKey *ecdsa.PublicKey) *owner.ID {
	return owner.NewIDFromPublicKey(pubKey)
}

func getPrivateKey(key *C.char) *keys.PrivateKey {
	return &keys.PrivateKey{PrivateKey: *getECDSAPrivKey(key)}
}

func getECDSAPrivKey(key *C.char) *ecdsa.PrivateKey {
	keyStr := C.GoString(key)
	bytes, err := hex.DecodeString(keyStr)
	die(err)
	k := new(big.Int)
	k.SetBytes(bytes)
	priv := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	priv.PublicKey.Curve = curve
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())
	return priv
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func getObjectIDFromV2(objectID *C.char) (*oid.ID, error) {
	id := new(oid.ID)
	err := id.UnmarshalJSON([]byte(C.GoString(objectID)))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal object id")
	}
	return id, nil
}

func getSessionTokenFromV2(sessionToken *C.char) (*session.Token, error) {
	token := new(session.Token)
	err := token.Unmarshal([]byte(C.GoString(sessionToken)))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal session token")
	}
	return token, nil
}

func getBearerTokenFromV2(bearerToken *C.char) (*token.BearerToken, error) {
	token := new(token.BearerToken)
	err := token.Unmarshal([]byte(C.GoString(bearerToken)))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal bearer token")
	}
	return token, nil
}

//func getTrustsFromV2(trust *C.char) (*[]reputation.Trust, error) {
//	t := new(reputation.Trust)
//	err := t.UnmarshalJSON([]byte(C.GoString(trust)))
//	if err != nil {
//		return nil, fmt.Errorf("could not unmarshal reputation trust")
//	}
//	return t, nil
//}

func getPeerToPeerTrustFromV2(p2pTrust *C.char) (*reputation.PeerToPeerTrust, error) {
	t := new(reputation.PeerToPeerTrust)
	err := t.UnmarshalJSON([]byte(C.GoString(p2pTrust)))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal peer to peer reputation trust")
	}
	return t, nil
}

//func getEpoch(epoch *C.char) uint64 {
//	return uint64(epoch)
//}

//func getIteration(iteration *C.char) uint32 {
//	return uint32(iteration)
//}

//func getSessionExpirationFromV2(expiration *C.ulong) uint64 {
//	return uint64(expiration)
//}

func getFiltersFromV2(filters *C.char) (*object.SearchFilters, error) {
	sfs := new(object.SearchFilters)
	err := sfs.UnmarshalJSON([]byte(C.GoString(filters)))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal search filters")
	}
	return sfs, nil
}

func getTableFromV2(table *C.char) (*eacl.Table, error) {
	tab := new(eacl.Table)
	err := tab.Unmarshal([]byte(C.GoString(table)))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal table")
	}
	return tab, nil
}

//func getAnnouncementsFromV2(announcement *C.char) []container.UsedSpaceAnnouncement {
//	c := new(container.UsedSpaceAnnouncement)
//	c.Unmarshal(C.GoString(announcement))
//}

func getBalanceRequestToSigned(req *v2accounting.BalanceRequest) *v2accounting.BalanceRequest {
	pr := &v2accounting.BalanceRequest{}
	m := pr.ToGRPCMessage().(proto.Message)
	json, err := message.MarshalJSON(req)
	die(err)
	println("balance json:")
	println(json)
	err = protojson.Unmarshal(json, m)
	die(err)

	println("unmarshalled:")
	println(m)
	err = pr.FromGRPCMessage(m)
	die(err)
	return pr
}

func getRequestToSigned(req *v2container.PutRequest) *v2container.PutRequest {
	pr := &v2container.PutRequest{}
	m := pr.ToGRPCMessage().(proto.Message)
	json, err := message.MarshalJSON(req)
	die(err)
	err = protojson.Unmarshal(json, m)
	die(err)

	err = pr.FromGRPCMessage(m)
	die(err)
	return pr
}

func prepareMetaHeaderBalancePut(req *v2accounting.BalanceRequest) {
	meta := req.GetMetaHeader()
	if meta == nil {
		meta = new(v2session.RequestMetaHeader)
		req.SetMetaHeader(meta)
	}
	if meta.GetTTL() == 0 {
		meta.SetTTL(2)
	}
	if meta.GetVersion() == nil {
		meta.SetVersion(version.Current().ToV2())
	}
	meta.SetNetworkMagic(12345)
}

func prepareMetaHeader(req *v2container.PutRequest) {
	meta := req.GetMetaHeader()
	if meta == nil {
		meta = new(v2session.RequestMetaHeader)
		req.SetMetaHeader(meta)
	}
	if meta.GetTTL() == 0 {
		meta.SetTTL(2)
	}
	if meta.GetVersion() == nil {
		meta.SetVersion(version.Current().ToV2())
	}
	meta.SetNetworkMagic(12345)
}

//endregion container old
//region helper

func getPubKey(publicKey *C.char) ecdsa.PublicKey {
	rawPub, _ := hex.DecodeString(C.GoString(publicKey))
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), rawPub)
	return ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
}

func getOwnerID(key ecdsa.PublicKey) owner.ID {
	return *owner.NewIDFromPublicKey(&key)
}

func getOwnerIDFromC(publicKey *C.char) owner.ID {
	return getOwnerID(getPubKey(publicKey))
}

//endregion helper
//region client

type NeoFSClient struct {
	mu     sync.RWMutex
	client *neofsCli.Client
}

type NeoFSClients struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]*NeoFSClient
}

func initClients(id uuid.UUID, newClient *neofsCli.Client) {
	neofsClients = &NeoFSClients{sync.RWMutex{}, map[uuid.UUID]*NeoFSClient{id: {sync.RWMutex{}, newClient}}}
}

func (clients *NeoFSClients) put(id uuid.UUID, newClient *neofsCli.Client) {
	clients.mu.Lock()
	clients.clients[id] = &NeoFSClient{
		mu:     sync.RWMutex{},
		client: newClient,
	}
	clients.mu.Unlock()
}

func (clients *NeoFSClients) delete(id uuid.UUID) {
	clients.mu.Lock()
	delete(clients.clients, id)
	clients.mu.Unlock()
}

var neofsClients *NeoFSClients

func getClient(clientID *C.char) (*NeoFSClient, error) {
	if neofsClients == nil {
		return nil, fmt.Errorf("no clients present")
	}
	cliID, err := uuid.Parse(C.GoString(clientID))
	if err != nil {
		return nil, fmt.Errorf("could not parse provided client id")
	}
	neofsClients.mu.RLock()
	cli := neofsClients.clients[cliID]
	if cli == nil {
		return nil, fmt.Errorf("no client present with id %v", C.GoString(clientID))
	}
	neofsClients.mu.RUnlock()
	return cli, nil
}

//export CreateClient
func CreateClient(key *C.char, neofsEndpoint *C.char) C.response {
	privateKey := getPrivateKey(key)
	endpoint := C.GoString(neofsEndpoint)
	newClient, err := neofsCli.New(
		neofsCli.WithDefaultPrivateKey(&privateKey.PrivateKey),
		neofsCli.WithURIAddress(endpoint, nil),
		neofsCli.WithNeoFSErrorParsing(),
	)
	if err != nil {
		return cResponseError(fmt.Errorf("cannot create neofs client: %w", err).Error())
	}
	u, err := uuid.NewUUID()
	if err != nil {
		return cResponseError("cannot create uuid")
	}

	if neofsClients == nil {
		initClients(u, newClient)
	} else {
		neofsClients.put(u, newClient)
	}
	return cResponseString("Client", u.String())
}

//endregion client

//region C.response

func cResponseError(errorMsg string) C.response {
	return C.response{C.CString("Error"), C.CString(errorMsg)}
}

func cResponseErrorStatus() C.response {
	return cResponseError("result status not successful")
}

func cResponseErrorClient() C.response {
	return cResponseError("could not get client")
}

func cResponse(responseType string, value []byte) C.response {
	return C.response{C.CString(responseType), (*C.char)(C.CBytes(value))}
}

func cResponseString(responseType string, value string) C.response {
	return C.response{C.CString(responseType), C.CString(value)}
}

//endregion C.response
