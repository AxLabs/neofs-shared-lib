package main

import "C"
import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	v2accounting "github.com/nspcc-dev/neofs-api-go/v2/accounting"
	v2container "github.com/nspcc-dev/neofs-api-go/v2/container"
	"github.com/nspcc-dev/neofs-api-go/v2/refs"
	"github.com/nspcc-dev/neofs-api-go/v2/rpc/message"
	v2session "github.com/nspcc-dev/neofs-api-go/v2/session"
	"github.com/nspcc-dev/neofs-api-go/v2/signature"
	crypto "github.com/nspcc-dev/neofs-crypto"
	"github.com/nspcc-dev/neofs-sdk-go/acl"
	neofsCli "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	cid "github.com/nspcc-dev/neofs-sdk-go/container/id"
	"github.com/nspcc-dev/neofs-sdk-go/eacl"
	"github.com/nspcc-dev/neofs-sdk-go/netmap"
	"github.com/nspcc-dev/neofs-sdk-go/object"
	oid "github.com/nspcc-dev/neofs-sdk-go/object/id"
	"github.com/nspcc-dev/neofs-sdk-go/owner"
	"github.com/nspcc-dev/neofs-sdk-go/reputation"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	"github.com/nspcc-dev/neofs-sdk-go/token"
	sigutil "github.com/nspcc-dev/neofs-sdk-go/util/signature"
	"github.com/nspcc-dev/neofs-sdk-go/version"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"math/big"
	"unsafe"
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

/*
----Accounting----
Balance
*/

/*
----Container----
Put
Get
Delete
List
SetExtendedACL
GetExtendedACL
AnnounceUsedSpace
*/

//export PutContainer
func PutContainer(neofsEndpoint *C.char, key *C.char, v2Container *C.char) *C.char {
	//TESTNET := "grpcs://st01.testnet.fs.neo.org:8082"
	privateKey := getPrivateKey(key)
	ownerAcc := wallet.NewAccountFromPrivateKey(privateKey)
	fsCli := getClient(key, neofsEndpoint)

	ctx := context.Background()

	// Parse the container
	cnr := getContainerFromV2(v2Container)

	// Overwrites potential set container version and owner id
	cnr.SetVersion(version.Current())
	cnr.SetOwnerID(getOwnerIDFromAccount(ownerAcc))

	// The following are expected to be set within the provided container parameter
	//  - placement policy
	//  - permissions
	//  - attributes

	var prmContainerPut neofsCli.PrmContainerPut
	prmContainerPut.SetContainer(*cnr)

	cnrResponse, err := fsCli.ContainerPut(ctx, prmContainerPut)
	if err != nil {
		panic(err)
	}

	containerID := cnrResponse.ID().String()
	cstr := C.CString(containerID)
	return cstr
}

//export GetContainer
func GetContainer(neofsEndpoint *C.char, key *C.char, containerID *C.char) unsafe.Pointer {
	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	// Parse the container
	id := getContainerIDFromV2(containerID)

	var prmContainerGet neofsCli.PrmContainerGet
	prmContainerGet.SetContainer(*id)

	cnrResponse, err := fsCli.ContainerGet(ctx, prmContainerGet)
	if err != nil {
		panic(err)
	}

	if !apistatus.IsSuccessful(cnrResponse.Status()) {
		return nil
	}
	containerJson, _ := cnrResponse.Container().MarshalJSON()
	return C.CBytes(containerJson)
}

//export DeleteContainer
func DeleteContainer(neofsEndpoint *C.char, key *C.char, containerID *C.char) *C.char {
	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	// Parse the container
	id := getContainerIDFromV2(containerID)

	var prmContainerGet neofsCli.PrmContainerDelete
	prmContainerGet.SetContainer(*id)
	//prmContainerGet.SetSessionToken()

	cnrResponse, err := fsCli.ContainerDelete(ctx, prmContainerGet)
	if err != nil {
		panic(err)
	}

	// Handle unsuccessful request
	if !apistatus.IsSuccessful(cnrResponse.Status()) {
		return nil
	}
	return nil
}

//export ListContainer
func ListContainer(neofsEndpoint *C.char, key *C.char, ownerPubKey *C.char) *C.char {
	//fsCli := getClient(key, neofsEndpoint)
	//ctx := context.Background()
	//
	//privateKey := getPrivateKey(key)
	//ownerAcc := wallet.NewAccountFromPrivateKey(privateKey)
	//
	//getOwnerIDFromPublicKey(ownerPubKey)
	//
	//// Parse the container
	//id := getContainerIDFromV2(containerID)
	//
	//var prmContainerList neofsCli.PrmContainerList
	//prmContainerList.SetAccount(ownerId)
	//prmContainerGet.SetContainer(*id)
	//
	//cnrResponse, err := fsCli.ContainerList(ctx, prmContainerGet)
	//if err != nil {
	//	panic(err)
	//}
	//
	//if !apistatus.IsSuccessful(cnrResponse.Status()) {
	//	return nil
	//}
	//containerJson, _ := cnrResponse.Container().MarshalJSON()
	//return C.CString(containerJson)
	return C.CString("")
}

//export SetExtendedACL
func SetExtendedACL(table *C.char, neofsEndpoint *C.char, key *C.char) *C.char {
	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	// Parse the table
	tab := getTable(table)

	var prmContainerSetEACL neofsCli.PrmContainerSetEACL
	prmContainerSetEACL.SetTable(*tab)

	cnrResponse, err := fsCli.ContainerSetEACL(ctx, prmContainerSetEACL)
	if err != nil {
		panic(err)
	}

	// Handle unsuccessful request
	if !apistatus.IsSuccessful(cnrResponse.Status()) {
		return nil
	}
	return nil
}

//export GetExtendedACL
func GetExtendedACL(containerID *C.char, neofsEndpoint *C.char, key *C.char) unsafe.Pointer {
	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	// Parse the container
	id := getContainerIDFromV2(containerID)

	var prmContainerEACL neofsCli.PrmContainerEACL
	prmContainerEACL.SetContainer(*id)

	cnrResponse, err := fsCli.ContainerEACL(ctx, prmContainerEACL)
	if err != nil {
		panic(err)
	}

	if !apistatus.IsSuccessful(cnrResponse.Status()) {
		return nil
	}
	containerJson, _ := cnrResponse.Table().MarshalJSON()
	return C.CBytes(containerJson)
}

////export AnnounceUsedSpace
//func AnnounceUsedSpace(announcements *C.char, neofsEndpoint *C.char, key *C.char) *C.char {
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	// Parse the container
//	ann := getAnnouncementsFromV2(announcements)
//
//	var prmContainerAnnounceSpace neofsCli.PrmAnnounceSpace
//	prmContainerAnnounceSpace.SetValues(ann)
//
//	cnrResponse, err := fsCli.ContainerAnnounceUsedSpace(ctx, prmContainerAnnounceSpace)
//	if err != nil {
//		panic(err)
//	}
//
//	if !apistatus.IsSuccessful(cnrResponse.Status()) {
//		return nil
//	}
//	containerJson, _ := cnrResponse.Container().MarshalJSON()
//	return C.CString(containerJson)
//}

/*
----Netmap----
LocalNodeInfo -> Could not find method
NetworkInfo
EndpointInfo
*/

//export NetworkInfo
func NetworkInfo(neofsEndpoint *C.char, key *C.char) unsafe.Pointer {
	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	var prmNetworkInfo neofsCli.PrmNetworkInfo
	response, err := fsCli.NetworkInfo(ctx, prmNetworkInfo)
	if err != nil {
		panic(err)
	}

	if !apistatus.IsSuccessful(response.Status()) {
		return nil
	}
	networkInfo, _ := response.Info().MarshalJSON()
	return C.CBytes(networkInfo)
}

//export EndpointInfo
func EndpointInfo(neofsEndpoint *C.char, key *C.char) unsafe.Pointer {
	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	var prmEndpointInfo neofsCli.PrmEndpointInfo
	response, err := fsCli.EndpointInfo(ctx, prmEndpointInfo)
	if err != nil {
		panic(err)
	}

	if !apistatus.IsSuccessful(response.Status()) {
		return nil
	}
	nodeInfo, _ := response.NodeInfo().MarshalJSON()
	return C.CBytes(nodeInfo)
}

func NetworkLatestVersion(neofsEndpoint *C.char, key *C.char) unsafe.Pointer {
	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	var prmEndpointInfo neofsCli.PrmEndpointInfo
	response, err := fsCli.EndpointInfo(ctx, prmEndpointInfo)
	if err != nil {
		panic(err)
	}

	if !apistatus.IsSuccessful(response.Status()) {
		return nil
	}
	latestVersion, _ := response.LatestVersion().MarshalJSON()
	return C.CBytes(latestVersion)
}

/*
----Object----
Get
Put
Delete
Head
Search
GetRange
GetRangeHash
*/

////export GetObjectInit
//func GetObjectInit(containerID *C.char, neofsEndpoint *C.char, key *C.char) *C.char {
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	// Parse the container
//	id := getContainerIDFromV2(containerID)
//
//	var prmObjectGet neofsCli.PrmObjectGet
//	prmObjectGet.FromContainer()
//
//	response, err := fsCli.ObjectGetInit(ctx, prmObjectGet)
//	if err != nil {
//		panic(err)
//	}
//
//	response.Read()
//	return C.CString(containerJson)
//}

////export ReadObject
//func ReadObject(objectReader *C.char) {
//
//}

////export PutObject
//func PutObject(neofsEndpoint *C.char, key *C.char) *C.char {
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	var prmObjectPutInit neofsCli.PrmObjectPutInit
//	response, err := fsCli.ObjectPutInit(ctx, prmObjectPutInit)
//	if err != nil {
//		panic(err)
//	}
//
//	response.WritePayloadChunk()
//	//return C.CString(containerJson)
//}

//export DeleteObject
func DeleteObject(containerID *C.char, objectID *C.char, sessionToken *C.char, bearerToken *C.char, neofsEndpoint *C.char,
	key *C.char) unsafe.Pointer {

	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	cid := getContainerIDFromV2(containerID)
	oidV2 := getObjectIDFromV2(objectID)
	st := getSessionTokenFromV2(sessionToken)
	bt := getBearerTokenFromV2(bearerToken)

	var prmObjectDelete neofsCli.PrmObjectDelete
	prmObjectDelete.FromContainer(*cid)
	prmObjectDelete.ByID(*oidV2)
	prmObjectDelete.WithinSession(*st)
	prmObjectDelete.WithBearerToken(*bt)

	cnrResponse, err := fsCli.ObjectDelete(ctx, prmObjectDelete)
	if err != nil {
		panic(err)
	}

	if !apistatus.IsSuccessful(cnrResponse.Status()) {
		return nil
	}
	dst := new(oid.ID)
	tombstoneRead := cnrResponse.ReadTombstoneID(dst)
	if !tombstoneRead {
		panic("Could not read object's tombstone.")
	}
	d, _ := dst.MarshalJSON()
	return C.CBytes(d)
}

//export GetObjectHead
func GetObjectHead(containerID *C.char, objectID *C.char, sessionToken *C.char, bearerToken *C.char, neofsEndpoint *C.char,
	key *C.char) unsafe.Pointer {

	fsCli := getClient(key, neofsEndpoint)
	ctx := context.Background()

	cid := getContainerIDFromV2(containerID)
	oid := getObjectIDFromV2(objectID)
	st := getSessionTokenFromV2(sessionToken)
	bt := getBearerTokenFromV2(bearerToken)

	var prmObjectHead neofsCli.PrmObjectHead
	prmObjectHead.FromContainer(*cid)
	prmObjectHead.ByID(*oid)
	prmObjectHead.WithinSession(*st)
	prmObjectHead.WithBearerToken(*bt)

	response, err := fsCli.ObjectHead(ctx, prmObjectHead)
	if err != nil {
		panic(err)
	}

	if !apistatus.IsSuccessful(response.Status()) {
		return nil
	}
	obj := new(object.Object)
	response.ReadHeader(obj)
	objectWithReadHeader, _ := obj.MarshalJSON()
	return C.CBytes(objectWithReadHeader)
}

////export SearchObject s?
//func SearchObject(containerID *C.char, sessionToken *C.char, bearerToken *C.char, filters *C.char, neofsEndpoint *C.char,
//	key *C.char) *C.char {
//
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	cid := getContainerIDFromV2(containerID)
//	st := getSessionTokenFromV2(sessionToken)
//	bt := getBearerTokenFromV2(bearerToken)
//	sfs := getFiltersFromV2(filters)
//
//	var prmObjectSearch neofsCli.PrmObjectSearch
//	prmObjectSearch.InContainer(*cid)
//	prmObjectSearch.WithinSession(*st)
//	prmObjectSearch.WithBearerToken(*bt)
//	prmObjectSearch.SetFilters(*sfs)
//
//	response, err := fsCli.ObjectSearchInit(ctx, prmObjectSearch)
//	if err != nil {
//		panic(err)
//	}
//
//	read, b := response.Read()
//	foundObject, _ := response.Close()
//
//	return C.CString(containerJson)
//}

////export GetRange
//func GetRange(containerID *C.char, objectID *C.char, sessionToken *C.char, bearerToken *C.char, length *C.char, offset *C.char,
//	neofsEndpoint *C.char, key *C.char) *C.char {
//
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	cid := getContainerIDFromV2(containerID)
//	oid := getObjectIDFromV2(objectID)
//	st := getSessionTokenFromV2(sessionToken)
//	bt := getBearerTokenFromV2(bearerToken)
//
//	var prmObjectRange neofsCli.PrmObjectRange
//	prmObjectRange.FromContainer(*cid)
//	prmObjectRange.ByID(*oid)
//	prmObjectRange.WithinSession(*st)
//	prmObjectRange.WithBearerToken(*bt)
//	prmObjectRange.SetLength(length)
//	prmObjectRange.SetOffset(offset)
//
//	response, err := fsCli.ObjectRangeInit(ctx, prmObjectRange)
//	if err != nil {
//		panic(err)
//	}
//
//	response.Read()
//}

////export GetRangeHash
//func GetRangeHash(containerID *C.char, neofsEndpoint *C.char, key *C.char) *C.char {
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	// Parse the container
//	id := getContainerIDFromV2(containerID)
//
//	var prmContainerGet neofsCli.PrmContainerGet
//	prmContainerGet.SetContainer(*id)
//
//	cnrResponse, err := fsCli.(ctx, prmContainerGet)
//	if err != nil {
//		panic(err)
//	}
//
//	if !apistatus.IsSuccessful(cnrResponse.Status()) {
//		return nil
//	}
//	containerJson, _ := cnrResponse.Container().MarshalJSON()
//	return C.CString(containerJson)
//}

/*
----Reputation----
AnnounceLocalTrust
AnnounceIntermediateResult
*/
////export AnnounceLocalTrust
//func AnnounceLocalTrust(containerID *C.char, trust *C.char, neofsEndpoint *C.char, key *C.char) *C.char {
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	// Parse the container
//	id := getContainerIDFromV2(containerID)
//	getTrustFromV2(trust)
//
//	var prmAnnounceLocalTrust neofsCli.PrmAnnounceLocalTrust
//	prmAnnounceLocalTrust.SetValues()
//	prmAnnounceLocalTrust.SetEpoch()
//
//	cnrResponse, err := fsCli.AnnounceLocalTrust(ctx, prmContainerGet)
//	if err != nil {
//		panic(err)
//	}
//
//	if !apistatus.IsSuccessful(cnrResponse.Status()) {
//		return nil
//	}
//	containerJson, _ := cnrResponse.Container().MarshalJSON()
//	return C.CString(containerJson)
//}

////export AnnounceIntermediateResult
//func AnnounceIntermediateResult(p2pTrust *C.char, epoch *C.char, iteration *C.char, neofsEndpoint *C.char, key *C.char) *C.char {
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	// Parse the container
//	trust := getPeerToPeerTrust(p2pTrust)
//	ep := getEpoch(epoch)         // uint64
//	it := getIteration(iteration) // uint32
//
//	var prmAnnounceIntermediateTrust neofsCli.PrmAnnounceIntermediateTrust
//	prmAnnounceIntermediateTrust.SetCurrentValue(*trust)
//	prmAnnounceIntermediateTrust.SetEpoch(ep)
//	prmAnnounceIntermediateTrust.SetIteration(it)
//
//	response, err := fsCli.AnnounceIntermediateTrust(ctx, prmAnnounceIntermediateTrust)
//	if err != nil {
//		panic(err)
//	}
//
//	// Handle unsuccessful request
//	if !apistatus.IsSuccessful(response.Status()) {
//		return nil
//	}
//	return nil
//}

/*
----Session----
Create
*/
////export CreateSession
//func CreateSession(sessionExpiration *C.ulonglong, neofsEndpoint *C.char, key *C.char) unsafe.Pointer {
//	fsCli := getClient(key, neofsEndpoint)
//	ctx := context.Background()
//
//	exp := getSessionExpirationFromV2(sessionExpiration)
//
//	var prmSessionCreate neofsCli.PrmSessionCreate
//	prmSessionCreate.SetExp(exp)
//
//	cnrResponse, err := fsCli.SessionCreate(ctx, prmSessionCreate)
//	if err != nil {
//		panic(err)
//	}
//
//	if !apistatus.IsSuccessful(cnrResponse.Status()) {
//		return nil
//	}
//
//	sessionID := cnrResponse.ID()
//	//sessionPubKey := cnrResponse.PublicKey()
//	return C.CBytes(sessionID)
//}

func getClient(key *C.char, neofsEndpoint *C.char) *neofsCli.Client {
	privateKey := getPrivateKey(key)
	endpoint := C.GoString(neofsEndpoint)
	cli, err := neofsCli.New(
		neofsCli.WithDefaultPrivateKey(&privateKey.PrivateKey),
		neofsCli.WithURIAddress(endpoint, nil),
		neofsCli.WithNeoFSErrorParsing(),
	)
	if err != nil {
		panic(fmt.Errorf("can't create neofs client: %w", err))
	}
	return cli
}

//export PutContainerBasic
func PutContainerBasic(key *C.char) *C.char {
	TESTNET := "grpcs://st01.testnet.fs.neo.org:8082"
	// create client from parameter
	//ctx := context.TODO()
	ctx := context.Background()
	//walletCli, err := client.New(ctx, "http://seed1t4.neo.org:2332", client.Options{}) // get Neo endpoint from parameter
	//if err != nil {
	//	return fmt.Errorf("can't create wallet client: %w", err)
	//}

	privateKey := keys.PrivateKey{PrivateKey: *getECDSAPrivKey(key)}
	ownerAcc := wallet.NewAccountFromPrivateKey(&privateKey)
	fsCli, err := neofsCli.New(
		neofsCli.WithDefaultPrivateKey(&privateKey.PrivateKey),
		neofsCli.WithURIAddress(TESTNET, nil), // get NeoFS endpoint from parameter
		neofsCli.WithNeoFSErrorParsing(),
	)
	if err != nil {
		panic(fmt.Errorf("can't create neofs client: %w", err))
	}

	//	create container from parameter
	//	required:
	//	o	create placement policy
	//	x	access to private key
	//	o	set permissions
	//	o	neofs client

	ownerID := getOwnerIDFromAccount(ownerAcc)

	placementPolicy := netmap.NewPlacementPolicy() // get placement policy from string

	permissions := acl.PublicBasicRule
	//acl.ParseBasicACL(aclString) // get acl from string argument

	cnr := container.New(
		container.WithPolicy(placementPolicy),
		container.WithOwnerID(ownerID),
		container.WithCustomBasicACL(permissions),
	)

	//attributes := container.Attributes{} // get attributes from string argument
	//cnr.SetAttributes(attributes)

	var prmContainerPut neofsCli.PrmContainerPut
	prmContainerPut.SetContainer(*cnr)

	cnrResponse, err := fsCli.ContainerPut(ctx, prmContainerPut)
	if err != nil {
		panic(err)
	}

	containerID := cnrResponse.ID().String()
	cstr := C.CString(containerID)
	return cstr
}

func getOwnerIDFromAccount(acc *wallet.Account) *owner.ID {
	return owner.NewIDFromN3Account(acc)
}

func getOwnerIDFromPublicKey(pubKey *ecdsa.PublicKey) *owner.ID {
	return owner.NewIDFromPublicKey(pubKey)
}

//export SignServiceMessage
func SignServiceMessage(key *C.char, json *C.char) *C.char {
	getECDSAPrivKey(key)
	keyStr := C.GoString(key)
	jsonStr := C.GoString(json)

	bytes, err := hex.DecodeString(keyStr)
	print(err)

	k := new(big.Int)
	k.SetBytes(bytes)

	priv := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	priv.PublicKey.Curve = curve
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())

	//this print can be used to verify if we got the same parameters as in Java version
	fmt.Printf("X: %d, Y: %d\n", priv.PublicKey.X, priv.PublicKey.Y)

	wif, err := crypto.WIFEncode(priv)
	fmt.Printf("WIF: %s\n", wif)

	pr := &v2container.PutRequest{}

	m := pr.ToGRPCMessage().(proto.Message)
	err = protojson.Unmarshal([]byte(jsonStr), m)
	if err != nil {
		panic(err)
	}

	err = pr.FromGRPCMessage(m)
	if err != nil {
		panic(err)
	}

	err = signature.SignServiceMessage(priv, pr)
	if err != nil {
		panic(err)
	}

	jsonAfter, err := message.MarshalJSON(pr)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonAfter))

	cstr := C.CString(string(jsonAfter))

	return cstr
}

//func getPublicKey(key *C.char) *keys.PublicKey {
//	return &keys.PublicKey{PublicKey: *getECDSAPubKey()}
//}
//
//func getECDSAPubKey(key *C.char) *ecdsa.PublicKey {
//	keyStr := C.GoString(key)
//}

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

//export VerifyServiceMessage
func VerifyServiceMessage(msg interface{}) error {
	return signature.VerifyServiceMessage(msg)
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func getContainerFromV2(v2Container *C.char) *container.Container {
	sdkContainer := new(container.Container)
	str := C.GoString(v2Container)
	err := sdkContainer.UnmarshalJSON([]byte(str))
	if err != nil {
		panic("Could not unmarshal container.")
	}
	return sdkContainer
}

func getContainerIDFromV2(containerID *C.char) *cid.ID {
	id := new(cid.ID)
	err := id.UnmarshalJSON([]byte(C.GoString(containerID)))
	if err != nil {
		panic("Could not unmarshal container id.")
	}
	return id
}

func getObjectIDFromV2(objectID *C.char) *oid.ID {
	id := new(oid.ID)
	err := id.UnmarshalJSON([]byte(C.GoString(objectID)))
	if err != nil {
		panic("Could not unmarshal object id.")
	}
	return id
}

func getSessionTokenFromV2(sessionToken *C.char) *session.Token {
	token := new(session.Token)
	err := token.Unmarshal([]byte(C.GoString(sessionToken)))
	if err != nil {
		panic("Could not unmarshal session token.")
	}
	return token
}

func getBearerTokenFromV2(bearerToken *C.char) *token.BearerToken {
	token := new(token.BearerToken)
	err := token.Unmarshal([]byte(C.GoString(bearerToken)))
	if err != nil {
		panic("Could not unmarshal bearer token.")
	}
	return token
}

func getTrustFromV2(trust *C.char) *reputation.Trust {
	t := new(reputation.Trust)
	err := t.UnmarshalJSON([]byte(C.GoString(trust)))
	if err != nil {
		panic("Could not unmarshal reputation trust.")
	}
	return t
}

func getPeerToPeerTrust(p2pTrust *C.char) *reputation.PeerToPeerTrust {
	t := new(reputation.PeerToPeerTrust)
	err := t.UnmarshalJSON([]byte(C.GoString(p2pTrust)))
	if err != nil {
		panic("Could not unmarshal peer to peer reputation trust.")
	}
	return t
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

func getFiltersFromV2(filters *C.char) *object.SearchFilters {
	sfs := new(object.SearchFilters)
	err := sfs.UnmarshalJSON([]byte(C.GoString(filters)))
	if err != nil {
		panic("Could not unmarshal search filters.")
	}
	return sfs
}

func getTable(table *C.char) *eacl.Table {
	tab := new(eacl.Table)
	err := tab.Unmarshal([]byte(C.GoString(table)))
	if err != nil {
		panic("Could not unmarshal table.")
	}
	return tab
}

//func getAnnouncementsFromV2(announcement *C.char) []container.UsedSpaceAnnouncement {
//	c := new(container.UsedSpaceAnnouncement)
//	c.Unmarshal(C.GoString(announcement))
//}

//export GetBalanceRequest
func GetBalanceRequest(key *C.char, ownerAddress *C.char) *C.char {
	privKey := getECDSAPrivKey(key)
	ownerIDString := C.GoString(ownerAddress)
	println("owner id string:")
	println(ownerIDString)
	ownerID := new(refs.OwnerID)
	ownerID.SetValue([]byte(ownerIDString))
	println("owner id getvalue:")
	println(ownerID.GetValue())
	var body v2accounting.BalanceRequestBody
	body.SetOwnerID(ownerID)

	var req v2accounting.BalanceRequest
	req.SetBody(&body)
	var meta v2session.RequestMetaHeader
	req.SetMetaHeader(&meta)
	prepareMetaHeaderBalancePut(&req)

	pr := getBalanceRequestToSigned(&req)

	err := signature.SignServiceMessage(privKey, pr)
	die(err)

	jsonAfter, err := message.MarshalJSON(pr)
	die(err)

	cstr := C.CString(string(jsonAfter))

	return cstr
}

//export NewContainerPutRequest
func NewContainerPutRequest(key *C.char, v2Container *C.char) *C.char {
	privKey := getECDSAPrivKey(key)

	cnr := getContainerFromV2(v2Container)
	if cnr.Version() == nil {
		cnr.SetVersion(version.Current())
	}
	_, err := cnr.NonceUUID()
	if err != nil {
		rand, err := uuid.NewRandom()
		if err != nil {
			panic("can't create new random " + err.Error())
		}
		cnr.SetNonceUUID(rand)
	}
	if cnr.BasicACL() == 0 {
		cnr.SetBasicACL(acl.PrivateBasicRule)
	}

	// form request body
	reqBody := new(v2container.PutRequestBody)
	reqBody.SetContainer(cnr.ToV2())

	// sign cnr
	signWrapper := signature.StableMarshalerWrapper{SM: reqBody.GetContainer()}
	err = sigutil.SignDataWithHandler(privKey, signWrapper, func(key []byte, sig []byte) {
		containerSignature := new(refs.Signature)
		containerSignature.SetKey(key)
		containerSignature.SetSign(sig)
		reqBody.SetSignature(containerSignature)
	}, sigutil.SignWithRFC6979())
	die(err)

	// form meta header
	var meta v2session.RequestMetaHeader
	meta.SetSessionToken(cnr.SessionToken().ToV2())

	// form request
	var req v2container.PutRequest
	req.SetBody(reqBody)

	// Prepare Meta Header
	// TODO: Check meta header params and set them accordingly
	// 	i.e., ttl, version, network magic
	req.SetMetaHeader(&meta)

	prepareMetaHeader(&req)

	pr := getRequestToSigned(&req)

	err = signature.SignServiceMessage(privKey, pr)
	die(err)

	jsonAfter, err := message.MarshalJSON(pr)
	die(err)

	cstr := C.CString(string(jsonAfter))

	return cstr
}

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

func getMessageCChar(req message.Message) (*C.char, error) {
	jsonAfter, err := message.MarshalJSON(req)
	die(err)
	return C.CString(string(jsonAfter)), nil
}
