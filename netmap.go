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
	neofsCli "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"reflect"
)

/*
----Netmap----
LocalNodeInfo?
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
		return pointerResponseError("could not get endpoint info")
	}
	status := resEndpointInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	latestVersion := resEndpointInfo.LatestVersion()
	if latestVersion == nil {
		return pointerResponseError("could not get latest version of endpoint")
	}
	nodeInfo := resEndpointInfo.NodeInfo()
	if nodeInfo == nil {
		return pointerResponseError("could not get node info of endpoint")
	}
	nodeInfoJson, err := nodeInfo.MarshalJSON()
	if err != nil {
		return pointerResponseError(err.Error())
	}
	latestVersionJson, err := latestVersion.MarshalJSON()
	if err != nil {
		return pointerResponseError(err.Error())
	}

	resp := EndpointResponse{
		NodeInfo:      string(nodeInfoJson),
		LatestVersion: string(latestVersionJson),
	}
	jsonArray, err := json.Marshal(resp)
	if err != nil {
		return pointerResponseError("could not marshal json")
	}
	return pointerResponse(reflect.TypeOf("EndpointInfo"), jsonArray)
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
	json, err := info.Marshal()
	if err != nil {
		return pointerResponseError("could not marshal network info of endpoint")
	}
	return pointerResponse(reflect.TypeOf(info), json)
}
