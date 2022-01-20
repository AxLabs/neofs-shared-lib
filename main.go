package main

import "C"
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/nspcc-dev/neofs-api-go/rpc/message"
	"github.com/nspcc-dev/neofs-api-go/v2/container"
	signpkg "github.com/nspcc-dev/neofs-api-go/v2/signature"
	crypto "github.com/nspcc-dev/neofs-crypto"
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

	pr := &container.PutRequest{}

	m := pr.ToGRPCMessage().(proto.Message)
	err = protojson.Unmarshal([]byte(jsonFromJava), m)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	_ = pr.FromGRPCMessage(m)

	err = signpkg.SignServiceMessage(priv, pr)
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

	pr := &container.PutRequest{}

	m := pr.ToGRPCMessage().(proto.Message)
	err = protojson.Unmarshal([]byte(jsonStr), m)
	if err != nil {
		panic(err)
	}

	err = pr.FromGRPCMessage(m)
	if err != nil {
		panic(err)
	}

	err = signpkg.SignServiceMessage(priv, pr)
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

//export VerifyServiceMessage
func VerifyServiceMessage(msg interface{}) error {
	return signpkg.VerifyServiceMessage(msg)
}
