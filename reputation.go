package main

/*
#include <stdlib.h>

#ifndef RESPONSE_H
#define RESPONSE_H
#include "response.h"
#endif
*/
import "C"

/*
----Reputation----
AnnounceLocalTrust
AnnounceIntermediateResult
*/

//export AnnounceLocalTrust
func AnnounceLocalTrust(clientID *C.char, v2Trusts *C.char, epoch *C.int) {}

//func AnnounceLocalTrust(clientID *C.char, v2Trusts *C.char, epoch *C.int) C.response {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return cResponseErrorClient()
//	}
//	cli.mu.RLock()
//	ctx := context.Background()
//	trusts, err := getTrustsFromV2(v2Trusts)
//	if err != nil {
//		return cResponseError(err.Error())
//	}
//
//	var prmAnnounceLocalTrust neofsCli.PrmAnnounceLocalTrust
//	prmAnnounceLocalTrust.SetValues(*trusts)
//	prmAnnounceLocalTrust.SetEpoch(C.GoInt(epoch))
//
//	resAnnounceLocalTrust, err := cli.client.AnnounceLocalTrust(ctx, prmAnnounceLocalTrust)
//	cli.mu.RUnlock()
//	if err != nil {
//		return cResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(resAnnounceLocalTrust.Status()) {
//		return cResponseErrorStatus()
//	}
//	return cResponse("AnnounceLocalTrust", ) // handle methods without return value
//}

//export AnnounceIntermediateResult
func AnnounceIntermediateResult(clientID *C.char, v2P2PTrust *C.char, epoch *C.char, iteration *C.char) {
}

//func AnnounceIntermediateResult(clientID *C.char, v2P2PTrust *C.char, epoch *C.char, iteration *C.char) C.response {
//	cli, err := getClient(clientID)
//	if err != nil {
//		return cResponseErrorClient()
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
//		return cResponseError(err.Error())
//	}
//	if !apistatus.IsSuccessful(resAnnounceIntermediateTrust.Status()) {
//		return cResponseErrorStatus()
//	}
//	return cResponse("AnnounceIntermediateLocalTrust", ) // handle methods without return value
//}
