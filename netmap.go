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
	"context"
	neofsCli "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
)

/*
----Netmap----
LocalNodeInfo?
NetworkInfo
EndpointInfo
*/

//export GetEndpointNodeInfo
func GetEndpointNodeInfo(clientID *C.char) C.response {
	cli, err := getClient(clientID)
	if cli.client == nil {
		return cResponseError("no client found")
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmEndpointInfo neofsCli.PrmEndpointInfo
	resEndpointInfo, err := cli.client.EndpointInfo(ctx, prmEndpointInfo)
	cli.mu.RUnlock()
	if err != nil {
		return cResponseError("could not get endpoint info")
	}
	status := resEndpointInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return cResponseErrorStatus()
	}
	info := resEndpointInfo.NodeInfo()
	if info == nil {
		return cResponseError("could not get node info of the endpoint")
	}
	json, err := info.MarshalJSON()
	if err != nil {
		return cResponseError("could not marshal node info of the endpoint")
	}
	return cResponse("GetEndpointNodeInfo", json)
}

//export GetEndpointLatestVersion
func GetEndpointLatestVersion(clientID *C.char) C.response {
	cli, err := getClient(clientID)
	if err != nil {
		return cResponseError(err.Error())
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmEndpointInfo neofsCli.PrmEndpointInfo
	resEndpointInfo, err := cli.client.EndpointInfo(ctx, prmEndpointInfo)
	cli.mu.RUnlock()
	if err != nil {
		return cResponseError("could not get endpoint info")
	}
	status := resEndpointInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return cResponseErrorStatus()
	}
	latestVersion := resEndpointInfo.LatestVersion()
	if latestVersion == nil {
		return cResponseError("could not get latest version of endpoint")
	}
	json, err := latestVersion.MarshalJSON()
	if err != nil {
		return cResponseError("could not marshal latest version of endpoint")
	}
	return cResponse("GetEndpointLatestVersion", json) // TODO: Get the response type of latestVersion and include this in the return.
}

//export GetNetworkInfo
func GetNetworkInfo(clientID *C.char) C.response {
	cli, err := getClient(clientID)
	if err != nil {
		return cResponseError(err.Error())
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmNetworkInfo neofsCli.PrmNetworkInfo
	resNetworkInfo, err := cli.client.NetworkInfo(ctx, prmNetworkInfo)
	cli.mu.RUnlock()
	if err != nil {
		return cResponseError("could not get endpoint info")
	}
	status := resNetworkInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return cResponseErrorStatus()
	}
	info := resNetworkInfo.Info()
	if info == nil {
		return cResponseError("could not get network info of endpoint")
	}
	json, err := info.MarshalJSON()
	if err != nil {
		return cResponseError("could not marshal network info of endpoint")
	}
	//In order to allocate the memory in c use a more up-to-date implementation of the following:
	//resp := (*C.struct_response)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_response{}))))
	//resp.responseType = (*C.char)(C.CString("net"))
	//resp.value = (*C.char)(C.CBytes(json))
	//return resp
	return cResponse("GetNetworkInfo", json)
}
