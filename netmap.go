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
	"encoding/json"
	v2Netmap "github.com/nspcc-dev/neofs-api-go/v2/netmap"
	v2Refs "github.com/nspcc-dev/neofs-api-go/v2/refs"
	neofsCli "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"reflect"
)

/*
----Netmap----
NetworkInfo
EndpointInfo
*/

//export GetEndpoint
func GetEndpoint(clientID *C.char) C.pointerResponse {
	cli, err := getClient(clientID)
	if err != nil {
		return pointerResponseError(err.Error())
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmEndpointInfo neofsCli.PrmEndpointInfo
	resEndpointInfo, err := cli.client.EndpointInfo(ctx, prmEndpointInfo)
	cli.mu.RUnlock()
	if err != nil {
		return pointerResponseError(err.Error())
	}
	status := resEndpointInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	latestVersion := resEndpointInfo.LatestVersion()
	if latestVersion == nil {
		return pointerResponseError(err.Error())
	}
	nodeInfo := resEndpointInfo.NodeInfo()
	if nodeInfo == nil {
		return pointerResponseError(err.Error())
	}
	nodeInfoJson, err := nodeInfo.MarshalJSON()
	if err != nil {
		return pointerResponseError(err.Error())
	}
	var v2 v2Refs.Version
	latestVersion.WriteToV2(&v2)
	latestVersionJson, err := v2.MarshalJSON()
	if err != nil {
		return pointerResponseError(err.Error())
	}

	endpointResponse := EndpointResponse{
		NodeInfo:      string(nodeInfoJson),
		LatestVersion: string(latestVersionJson),
	}
	bytes, err := json.Marshal(endpointResponse)
	if err != nil {
		return pointerResponseError(err.Error())
	}
	return pointerResponse(reflect.TypeOf(endpointResponse), bytes)
}

type EndpointResponse struct {
	NodeInfo      string `json:"netmap.NodeInfo"`
	LatestVersion string `json:"version.Version"`
}

//export GetNetworkInfo
func GetNetworkInfo(clientID *C.char) C.pointerResponse {
	cli, err := getClient(clientID)
	if err != nil {
		return pointerResponseError(err.Error())
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmNetworkInfo neofsCli.PrmNetworkInfo
	resNetworkInfo, err := cli.client.NetworkInfo(ctx, prmNetworkInfo)
	cli.mu.RUnlock()
	if err != nil {
		return pointerResponseError("could not get endpoint info")
	}
	status := resNetworkInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	info := resNetworkInfo.Info()
	if info == nil {
		return pointerResponseError("could not get network info of endpoint")
	}
	var v2 v2Netmap.NetworkInfo
	info.WriteToV2(&v2)
	bytes := v2.StableMarshal(nil)

	if err != nil {
		return pointerResponseError("could not marshal network info of endpoint")
	}
	return pointerResponse(reflect.TypeOf(*info), bytes)
}
