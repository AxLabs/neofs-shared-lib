package session

///*
//----Session----
//Create
//*/
//
////export CreateSession
//func CreateSession(clientID *C.char, sessionExpiration *C.char) C.PointerResponse {
//	cli, err := GetClient(clientID)
//	if err != nil {
//		return PointerResponseClientError()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	var prmSessionCreate neofsclient.PrmSessionCreate
//	exp, err := strconv.ParseUint(C.GoString(sessionExpiration), 10, 64)
//	if err != nil {
//		return PointerResponseError("could not parse session expiration to uint64")
//	}
//	prmSessionCreate.SetExp(exp)
//
//	resSessionCreate, err := cli.client.SessionCreate(ctx, prmSessionCreate)
//	cli.mu.RUnlock()
//	if err != nil {
//		return PointerResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(resSessionCreate.Status()) {
//		return ResultStatusErrorResponsePointer()
//	}
//	sessionID := resSessionCreate.ID()
//	//sessionPublicKey := resSessionCreate.PublicKey()
//	return PointerResponse(reflect.TypeOf(sessionID), sessionID) // handle method with two return values
//}
//
////export CreateSessionPubKey
//func CreateSessionPubKey(clientID *C.char, sessionExpiration *C.char) C.PointerResponse {
//	cli, err := GetClient(clientID)
//	if err != nil {
//		return PointerResponseClientError()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	var prmSessionCreate neofsclient.PrmSessionCreate
//	exp, err := strconv.ParseUint(C.GoString(sessionExpiration), 10, 64)
//	if err != nil {
//		return PointerResponseError("could not parse session expiration to uint64")
//	}
//	prmSessionCreate.SetExp(exp)
//
//	resSessionCreate, err := cli.client.SessionCreate(ctx, prmSessionCreate)
//	cli.mu.RUnlock()
//	if err != nil {
//		return PointerResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(resSessionCreate.Status()) {
//		return ResultStatusErrorResponsePointer()
//	}
//	sessionPublicKey := resSessionCreate.PublicKey()
//	return PointerResponse(reflect.TypeOf(sessionPublicKey), sessionPublicKey) // handle method with two return values
//}
