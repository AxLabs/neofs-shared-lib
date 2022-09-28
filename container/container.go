package container

import (
	"context"
	"github.com/AxLabs/neofs-api-shared-lib/client"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	v2container "github.com/nspcc-dev/neofs-api-go/v2/container"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	cid "github.com/nspcc-dev/neofs-sdk-go/container/id"
	"reflect"
)

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

func PutContainer(neofsClient *client.NeoFSClient, cnr *container.Container) *response.StringResponse {
	ctx := context.TODO()

	var prmContainerPut neofsclient.PrmContainerPut
	prmContainerPut.SetContainer(*cnr)

	resContainerPut, err := neofsClient.LockAndGet().ContainerPut(ctx, prmContainerPut)
	neofsClient.Unlock()
	if err != nil {
		return response.StringError(err)
	}

	if !apistatus.IsSuccessful(resContainerPut.Status()) {
		return response.StringStatusResponse()
	}

	containerID := *resContainerPut.ID()
	return response.NewString(reflect.TypeOf(containerID), containerID.String())
}

func GetContainer(neofsClient *client.NeoFSClient, containerID *cid.ID) *response.PointerResponse {
	ctx := context.Background()

	var prmContainerGet neofsclient.PrmContainerGet
	prmContainerGet.SetContainer(*containerID)
	//prmContainerGet.WithXHeaders()

	resContainerGet, err := neofsClient.LockAndGet().ContainerGet(ctx, prmContainerGet)
	neofsClient.Unlock()

	if err != nil {
		return response.Error(err)
	}
	if !apistatus.IsSuccessful(resContainerGet.Status()) {
		return response.StatusResponse()
	}

	cnr := resContainerGet.Container()
	var v2 v2container.Container
	cnr.WriteToV2(&v2)
	if err != nil {
		return response.Error(err)
	}
	bytes := v2.StableMarshal(nil)
	return response.New(reflect.TypeOf(v2), bytes)
}

////export DeleteContainer
//func DeleteContainer(clientID *C.char, containerID *C.char) C.pointerResponse {
//	id, err := getContainerIDFromC(containerID)
//	if err != nil {
//		return response.PointerResponseError(err.Error())
//	}
//
//	neofsClient, err := client.GetClient(clientID)
//	if err != nil {
//		return response.PointerResponseClientError()
//	}
//	return deleteContainer(neofsClient, id, nil)
//}
//
//////export DeleteContainerWithinSession
////func DeleteContainerWithinSession(clientID *C.char, containerID *C.char, sessionToken *C.char) C.PointerResponse {
////	id, err := getContainerIDFromC(containerID)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////
////	tok, err := getSessionTokenFromC(sessionToken)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////
////	neofsClient, err := GetClient(clientID)
////	if err != nil {
////		return PointerResponseClientError()
////	}
////	neofsClient.mu.Lock()
////	return deleteContainer(neofsClient, id, tok)
////}
//
//func deleteContainer(neofsClient *client.NeoFSClient, containerID *cid.ID, sessionToken *session.Container) C.pointerResponse {
//	client := neofsClient.LockAndGet()
//	ctx := context.Background()
//
//	var prmContainerDelete neofsclient.PrmContainerDelete
//	prmContainerDelete.SetContainer(*containerID)
//	if sessionToken != nil {
//		prmContainerDelete.WithinSession(*sessionToken)
//	}
//	//prmContainerDelete.WithXHeaders()
//
//	resContainerDelete, err := client.ContainerDelete(ctx, prmContainerDelete)
//	neofsClient.Unlock()
//	if err != nil {
//		response.PointerResponseError(err.Error())
//	}
//
//	if !apistatus.IsSuccessful(resContainerDelete.Status()) {
//		return response.ResultStatusErrorResponsePointer()
//	}
//
//	return response.PointerResponseBoolean(true)
//}
//
////export ListContainer
//func ListContainer(clientID *C.char, ownerPubKey *C.char) C.pointerResponse {
//	ctx := context.Background()
//
//	var prmContainerList neofsclient.PrmContainerList
//	key, err := util.UserIDFromPublicKey(ownerPubKey)
//	if err != nil {
//		return response.PointerResponseError(err.Error())
//	}
//	prmContainerList.SetAccount(*key)
//	//prmContainerList.WithXHeaders()
//
//	neofsClient, err := client.GetClient(clientID)
//	if err != nil {
//		return response.PointerResponseClientError()
//	}
//	resContainerList, err := neofsClient.LockAndGet().ContainerList(ctx, prmContainerList)
//	neofsClient.Unlock()
//	if err != nil {
//		return response.PointerResponseError(err.Error())
//	}
//
//	if !apistatus.IsSuccessful(resContainerList.Status()) {
//		return response.ResultStatusErrorResponsePointer()
//	}
//
//	containerIDs := resContainerList.Containers()
//	ids := parseContainerIDs(containerIDs)
//	return response.PointerResponse(reflect.TypeOf(containerIDs), ids) // how return []cid.ID
//}
//
//func parseContainerIDs(containerIDList []cid.ID) []byte {
//	bytes := make([]byte, 0)
//	for _, id := range containerIDList {
//		bytes = append(bytes[:], []byte(id.EncodeToString())...)
//	}
//	return bytes
//}
//
//////export SetExtendedACL
////func SetExtendedACL(clientID *C.char, v2Table *C.char) C.PointerResponse {
////	cli, err := GetClient(clientID)
////	if err != nil {
////		return PointerResponseClientError()
////	}
////	cli.mu.RLock()
////	ctx := context.Background()
////	table, err := getTableFromV2(v2Table)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	var prmContainerSetEACL neofsclient.PrmContainerSetEACL
////	prmContainerSetEACL.SetTable(*table)
////
////	resContainerSetEACL, err := cli.client.ContainerSetEACL(ctx, prmContainerSetEACL)
////	cli.mu.RUnlock()
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	if !apistatus.IsSuccessful(resContainerSetEACL.Status()) {
////		return ResultStatusErrorResponsePointer()
////	}
////	boolean := []byte{1}
////	return PointerResponse(reflect.TypeOf(boolean), boolean)
////}
//
//////export GetExtendedACL
////func GetExtendedACL(clientID *C.char, v2ContainerID *C.char) C.PointerResponse {
////	cli, err := GetClient(clientID)
////	if err != nil {
////		return PointerResponseClientError()
////	}
////	cli.mu.RLock()
////	ctx := context.Background()
////	containerID, err := getV2ContainerIDFromC(v2ContainerID)
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	var prmContainerEACL neofsclient.PrmContainerEACL
////	prmContainerEACL.SetContainer(*containerID)
////
////	cnrResponse, err := cli.client.ContainerEACL(ctx, prmContainerEACL)
////	cli.mu.RUnlock()
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	if !apistatus.IsSuccessful(cnrResponse.Status()) {
////		return ResultStatusErrorResponsePointer()
////	}
////	table := cnrResponse.Table()
////	tableBytes, err := cnrResponse.Table().Marshal()
////	if err != nil {
////		return PointerResponseError("could not marshal eacl table")
////	}
////	return PointerResponse(reflect.TypeOf(table), tableBytes)
////}
//
//////export AnnounceUsedSpace
////func AnnounceUsedSpace(clientID *C.char, announcements *C.char) C.PointerResponse {
////	cli, err := GetClient(clientID)
////	if err != nil {
////		return PointerResponseClientError()
////	}
////	cli.mu.RLock()
////	ctx := context.Background()
////	ann := getAnnouncementsFromV2(announcements)
////
////	var prmContainerAnnounceSpace neofsclient.PrmAnnounceSpace
////	prmContainerAnnounceSpace.SetValues(ann)
////
////	resContainerAnnounceUsedSpace, err := cli.client.ContainerAnnounceUsedSpace(ctx, prmContainerAnnounceSpace)
////	cli.mu.RUnlock()
////	if err != nil {
////		return PointerResponseError(err.Error())
////	}
////	if !apistatus.IsSuccessful(resContainerAnnounceUsedSpace.Status()) {
////		return ResultStatusErrorResponsePointer()
////	}
////	boolean := []byte{1}
////	return PointerResponse(reflect.TypeOf(boolean), boolean)
////}
//
////endregion container
////region helper
//
//func getContainerFromC(v2Container *C.char) (*container.Container, error) {
//	v2cnr := new(v2container.Container)
//	err := v2cnr.UnmarshalJSON([]byte(C.GoString(v2Container)))
//	if err != nil {
//		return nil, err
//	}
//	//v2cnr.SetHomomorphicHashingState()
//
//	var cnr container.Container
//	err = cnr.ReadFromV2(*v2cnr)
//	if err != nil {
//		return nil, err
//	}
//	return &cnr, nil
//}
//
//func getContainerIDFromC(containerID *C.char) (*cid.ID, error) {
//	id := new(cid.ID)
//	err := id.DecodeString(C.GoString(containerID))
//	if err != nil {
//		return nil, err
//	}
//	return id, nil
//}
//
////func getSessionTokenFromC(sessionToken *C.char) (*session.Container, error) {
////	token := new(session.Container)
////
////	err := token.Unmarshal([]byte(C.GoString(sessionToken)))
////	if err != nil {
////		return nil, err
////	}
////	return token, nil
////}
//
////endregion helper
////region container old
//
//////export PutContainerBasic
////func PutContainerBasic(key *C.char) *C.char {
////	TESTNET := "grpcs://st01.testnet.fs.neo.org:8082"
////	// create client from parameter
////	//ctx := context.TODO()
////	ctx := context.Background()
////	//walletCli, err := client.New(ctx, "http://seed1t4.neo.org:2332", client.Options{}) // get Neo endpoint from parameter
////	//if err != nil {
////	//	return fmt.Errorf("can't create wallet client: %w", err)
////	//}
////
////	privateKey := keys.PrivateKey{PrivateKey: *GetECDSAPrivKey(key)}
////	ownerAcc := wallet.NewAccountFromPrivateKey(&privateKey)
////	fsCli, err := neofsclient.New(
////		neofsclient.WithDefaultPrivateKey(&privateKey.PrivateKey),
////		neofsclient.WithURIAddress(TESTNET, nil), // get NeoFS endpoint from parameter
////		neofsclient.WithNeoFSErrorParsing(),
////	)
////	if err != nil {
////		panic(fmt.Errorf("can't create neofs client: %w", err))
////	}
////
////	//	create container from parameter
////	//	required:
////	//	o	create placement policy
////	//	x	access to private key
////	//	o	set permissions
////	//	o	neofs client
////
////	ownerID := getOwnerIDFromAccount(ownerAcc)
////
////	placementPolicy := netmap.NewPlacementPolicy() // get placement policy from string
////
////	permissions := acl.PublicBasicRule
////	//acl.ParseBasicACL(aclString) // get acl from string argument
////
////	cnr := container.New(
////		container.WithPolicy(placementPolicy),
////		container.WithOwnerID(ownerID),
////		container.WithCustomBasicACL(permissions),
////	)
////
////	//attributes := container.Attributes{} // get attributes from string argument
////	//cnr.SetAttributes(attributes)
////
////	var prmContainerPut neofsclient.PrmContainerPut
////	prmContainerPut.SetContainer(*cnr)
////
////	cnrResponse, err := fsCli.ContainerPut(ctx, prmContainerPut)
////	if err != nil {
////		panic(err)
////	}
////
////	containerID := cnrResponse.ID().String()
////	cstr := C.CString(containerID)
////	return cstr
////}
//
//// old code
//////export NewContainerPutRequest
////func NewContainerPutRequest(key *C.char, v2Container *C.char) *C.char {
////	privKey := GetECDSAPrivKey(key)
////
////	cnr, err := getV2ContainerFromC(v2Container)
////	if err != nil {
////		panic("could not get container from v2")
////	}
////	if cnr.Version() == nil {
////		cnr.SetVersion(version.Current())
////	}
////	_, err = cnr.NonceUUID()
////	if err != nil {
////		rand, err := uuid.NewRandom()
////		if err != nil {
////			panic("can't create new random " + err.Error())
////		}
////		cnr.SetNonceUUID(rand)
////	}
////	if cnr.BasicACL() == 0 {
////		cnr.SetBasicACL(acl.PrivateBasicRule)
////	}
////
////	// form request body
////	reqBody := new(v2container.PutRequestBody)
////	reqBody.SetContainer(cnr.ToV2())
////
////	// sign cnr
////	signWrapper := signature.StableMarshalerWrapper{SM: reqBody.GetContainer()}
////	err = sigutil.SignDataWithHandler(privKey, signWrapper, func(key []byte, sig []byte) {
////		containerSignature := new(refs.Signature)
////		containerSignature.SetKey(key)
////		containerSignature.SetSign(sig)
////		reqBody.SetSignature(containerSignature)
////	}, sigutil.SignWithRFC6979())
////	die(err)
////
////	// form meta header
////	var meta v2session.RequestMetaHeader
////	meta.SetSessionToken(cnr.SessionToken().ToV2())
////
////	// form request
////	var req v2container.PutRequest
////	req.SetBody(reqBody)
////
////	// Prepare Meta Header
////	// TODO: Check meta header params and set them accordingly
////	// 	i.e., ttl, version, network magic
////	req.SetMetaHeader(&meta)
////
////	prepareMetaHeader(&req)
////
////	pr := getRequestToSigned(&req)
////
////	err = signature.SignServiceMessage(privKey, pr)
////	die(err)
////
////	jsonAfter, err := message.MarshalJSON(pr)
////	die(err)
////
////	cstr := C.CString(string(jsonAfter))
////
////	return cstr
////}
//
////endregion container old
