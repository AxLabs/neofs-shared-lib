package reputation

/*
----Reputation----
-AnnounceLocalTrust
-AnnounceIntermediateResult
*/

//
////export AnnounceLocalTrust
//func AnnounceLocalTrust(clientID *C.char, v2Trusts *C.char, epoch *C.int) C.Response {
//	cli, err := GetClient(clientID)
//	if err != nil {
//		return ResponseClientError()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	trusts, err := getTrustsFromV2(v2Trusts)
//	if err != nil {
//		return ResponseError(err.Error())
//	}
//
//	var prmAnnounceLocalTrust neofsCli.PrmAnnounceLocalTrust
//	prmAnnounceLocalTrust.SetValues(*trusts)
//	prmAnnounceLocalTrust.SetEpoch(C.GoInt(epoch))
//
//	resAnnounceLocalTrust, err := cli.client.AnnounceLocalTrust(ctx, prmAnnounceLocalTrust)
//	cli.mu.RUnlock()
//	if err != nil {
//		return ResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(resAnnounceLocalTrust.Status()) {
//		return ResultStatusErrorResponse()
//	}
//	boolean := []byte{1}
//	return Response(reflect.TypeOf(boolean), boolean)
//}

////export AnnounceIntermediateResult
//func AnnounceIntermediateResult(clientID *C.char, v2P2PTrust *C.char, epoch *C.char, iteration *C.char) C.Response {
//	cli, err := GetClient(clientID)
//	if err != nil {
//		return ResponseClientError()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	p2pTrust, err := getPeerToPeerTrustFromV2(v2P2PTrust)
//	ep := epoch     // uint64
//	it := iteration // uint32
//
//	var prmAnnounceIntermediateTrust neofsCli.PrmAnnounceIntermediateTrust
//	prmAnnounceIntermediateTrust.SetCurrentValue(*p2pTrust)
//	prmAnnounceIntermediateTrust.SetEpoch(ep)
//	prmAnnounceIntermediateTrust.SetIteration(it)
//
//	resAnnounceIntermediateTrust, err := cli.client.AnnounceIntermediateTrust(ctx, prmAnnounceIntermediateTrust)
//	cli.mu.RUnlock()
//	if err != nil {
//		return ResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(resAnnounceIntermediateTrust.Status()) {
//		return ResultStatusErrorResponse()
//	}
//	boolean := []byte{1}
//	return Response(reflect.TypeOf(boolean), boolean)
//}
