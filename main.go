package main

import "C"
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	v2accounting "github.com/nspcc-dev/neofs-api-go/v2/accounting"
	v2container "github.com/nspcc-dev/neofs-api-go/v2/container"
	"github.com/nspcc-dev/neofs-api-go/v2/refs"
	"github.com/nspcc-dev/neofs-api-go/v2/rpc/message"
	v2session "github.com/nspcc-dev/neofs-api-go/v2/session"
	"github.com/nspcc-dev/neofs-api-go/v2/signature"
	crypto "github.com/nspcc-dev/neofs-crypto"
	"github.com/nspcc-dev/neofs-sdk-go/acl"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	sigutil "github.com/nspcc-dev/neofs-sdk-go/util/signature"
	"github.com/nspcc-dev/neofs-sdk-go/version"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"math/big"
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

//export SignServiceMessage
func SignServiceMessage(key *C.char, json *C.char) *C.char {
	getPrivKey(key)
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

func getPrivKey(key *C.char) *ecdsa.PrivateKey {
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

func getContainerFromV2(v2Container *C.char) (*container.Container, error) {
	sdkContainer := new(container.Container)
	str := C.GoString(v2Container)
	err := sdkContainer.UnmarshalJSON([]byte(str))
	if err != nil {
		return nil, err
	}
	return sdkContainer, nil
}

//export GetBalance
func GetBalance(key *C.char, ownerAddress *C.char) *C.char {
	privKey := getPrivKey(key)
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
	privKey := getPrivKey(key)

	cnr, err := getContainerFromV2(v2Container)
	if cnr.Version() == nil {
		cnr.SetVersion(version.Current())
	}
	_, err = cnr.NonceUUID()
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
