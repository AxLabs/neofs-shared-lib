package main

import "C"
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	v2container "github.com/nspcc-dev/neofs-api-go/v2/container"
	"github.com/nspcc-dev/neofs-api-go/v2/refs"
	"github.com/nspcc-dev/neofs-api-go/v2/rpc/message"
	v2session "github.com/nspcc-dev/neofs-api-go/v2/session"
	v2signature "github.com/nspcc-dev/neofs-api-go/v2/signature"
	crypto "github.com/nspcc-dev/neofs-crypto"
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
	print(err)

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

	err = v2signature.SignServiceMessage(priv, pr)
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

	err = v2signature.SignServiceMessage(priv, pr)
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
	return v2signature.VerifyServiceMessage(msg)
}

type (
	//CallOption func(*callOptions)

	callOptions struct {
		version *version.Version
		//xHeaders []*session.XHeader
		ttl uint32
		//epoch    uint64
		//key      *ecdsa.PrivateKey
		//session  *session.Token
		//bearer   *token.BearerToken
	}
)

func defaultCallOptions() *callOptions {
	return &callOptions{
		version: version.Current(),
		ttl:     2,
	}
}

func v2MetaHeaderFromOpts(options *callOptions) *v2session.RequestMetaHeader {
	meta := new(v2session.RequestMetaHeader)
	meta.SetVersion(options.version.ToV2())
	meta.SetTTL(options.ttl)
	//meta.SetEpoch(options.epoch)

	//xhdrs := make([]*v2session.XHeader, len(options.xHeaders))
	//for i := range options.xHeaders {
	//	xhdrs[i] = options.xHeaders[i].ToV2()
	//}

	//meta.SetXHeaders(xhdrs)

	//if options.bearer != nil {
	//	meta.SetBearerToken(options.bearer.ToV2())
	//}

	//meta.SetSessionToken(options.session.ToV2())

	return meta
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func getContainerFromV2(v2Container *C.char) (*container.Container, error) {
	sdkContainer := new(container.Container)
	containerBytes, err := hex.DecodeString(C.GoString(v2Container))
	if err != nil {
		return nil, err
	}
	err = sdkContainer.Unmarshal(containerBytes)
	if err != nil {
		return nil, err
	}
	return sdkContainer, nil
}

func NewContainerPutRequest(key *C.char, v2Container *C.char) (*C.char, error) {
	privKey := getPrivKey(key)

	cnr, err := getContainerFromV2(v2Container)
	if err != nil {
		return nil, err
	}

	// form request body
	reqBody := new(v2container.PutRequestBody)
	reqBody.SetContainer(cnr.ToV2())

	// sign cnr
	signWrapper := v2signature.StableMarshalerWrapper{SM: reqBody.GetContainer()}
	err = sigutil.SignDataWithHandler(privKey, signWrapper, func(key []byte, sig []byte) {
		containerSignature := new(refs.Signature)
		containerSignature.SetKey(key)
		containerSignature.SetSign(sig)
		reqBody.SetSignature(containerSignature)
	}, sigutil.SignWithRFC6979())
	if err != nil {
		return nil, err
	}

	// form meta header
	var meta v2session.RequestMetaHeader
	meta.SetSessionToken(cnr.SessionToken().ToV2())

	// form request
	var req v2container.PutRequest
	req.SetBody(reqBody)
	req.SetMetaHeader(&meta)
	return getMessageCChar(&req)
}

func getMessageCChar(req message.Message) (*C.char, error) {
	jsonAfter, err := message.MarshalJSON(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonAfter))

	return C.CString(string(jsonAfter)), nil
}
