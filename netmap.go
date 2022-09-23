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
	v2netmap "github.com/nspcc-dev/neofs-api-go/v2/netmap"
	v2refs "github.com/nspcc-dev/neofs-api-go/v2/refs"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"reflect"
)

/*
----Netmap----
NetworkInfo
EndpointInfo
NetMapSnapshot (only exists >v1.0.0-rc.6)
*/

//export GetEndpoint
func GetEndpoint(clientID *C.char) C.pointerResponse {
	ctx := context.Background()
	var prmEndpointInfo neofsclient.PrmEndpointInfo

	neofsClient, err := getClient(clientID)
	if err != nil {
		return pointerResponseError(err.Error())
	}
	neofsClient.mu.Lock()
	resEndpointInfo, err := neofsClient.client.EndpointInfo(ctx, prmEndpointInfo)
	neofsClient.mu.Unlock()
	if err != nil {
		return pointerResponseError(err.Error())
	}

	// Todo: Return specific status instead of default unsuccessful status.
	if !apistatus.IsSuccessful(resEndpointInfo.Status()) {
		return resultStatusErrorResponsePointer()
	}

	bytes, err := buildEndpointResponse(resEndpointInfo)
	if err != nil {
		return pointerResponseError(err.Error())
	}
	return pointerResponse(reflect.TypeOf(neofsclient.ResEndpointInfo{}), bytes)
}

func buildEndpointResponse(resEndpointInfo *neofsclient.ResEndpointInfo) ([]byte, error) {
	latestVersion := resEndpointInfo.LatestVersion()
	nodeInfo := resEndpointInfo.NodeInfo()
	nodeInfoJson, err := nodeInfo.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var v2 v2refs.Version
	latestVersion.WriteToV2(&v2)
	latestVersionJson, err := v2.MarshalJSON()
	if err != nil {
		return nil, err
	}
	endpointResponse := EndpointResponse{
		NodeInfo:      string(nodeInfoJson),
		LatestVersion: string(latestVersionJson),
	}
	bytes, err := json.Marshal(endpointResponse)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type EndpointResponse struct {
	NodeInfo      string `json:"netmap.NodeInfo"`
	LatestVersion string `json:"version.Version"`
}

//export GetNetworkInfo
func GetNetworkInfo(clientID *C.char) C.pointerResponse {
	ctx := context.Background()
	var prmNetworkInfo neofsclient.PrmNetworkInfo
	//prmNetworkInfo.WithXHeaders()

	neofsClient, err := getClient(clientID)
	if err != nil {
		return pointerResponseError(err.Error())
	}
	neofsClient.mu.Lock()
	resNetworkInfo, err := neofsClient.client.NetworkInfo(ctx, prmNetworkInfo)
	neofsClient.mu.Unlock()
	if err != nil {
		return pointerResponseError(err.Error())
	}

	status := resNetworkInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	info := resNetworkInfo.Info()
	if info == nil {
		return pointerResponseError("could not get network info of endpoint")
	}

	var v2 v2netmap.NetworkInfo
	info.WriteToV2(&v2)
	bytes := v2.StableMarshal(nil)

	return pointerResponse(reflect.TypeOf(*info), bytes)
}
