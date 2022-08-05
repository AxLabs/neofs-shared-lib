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
	"fmt"
	neofsCli "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/netmap"
	"github.com/nspcc-dev/neofs-sdk-go/version"
	"reflect"
)

/*
----Netmap----
LocalNodeInfo?
NetworkInfo
EndpointInfo
*/

//export GetEndpoint
func GetEndpoint(clientID *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmEndpointInfo neofsCli.PrmEndpointInfo
	resEndpointInfo, err := cli.client.EndpointInfo(ctx, prmEndpointInfo)
	cli.mu.RUnlock()
	if err != nil {
		return errorResponsePointer("could not get endpoint info")
	}
	status := resEndpointInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	latestVersion := resEndpointInfo.LatestVersion()
	if latestVersion == nil {
		return errorResponsePointer("could not get latest version of endpoint")
	}
	nodeInfo := resEndpointInfo.NodeInfo()
	if nodeInfo == nil {
		return errorResponsePointer("could not get node info of endpoint")
	}
	//Todo: Find a way to return latest version and node info in one method
	//endpoint, err := newEndpoint(latestVersion, nodeInfo)
	//if err != nil {
	//	return errorResponse(err.Error())
	//}
	//resp := endpointResponse(*endpoint)
	marshal, _ := nodeInfo.Marshal()
	return newResponsePointer(reflect.TypeOf(nodeInfo), marshal)
}

//export GetEndpointNodeInfo
func GetEndpointNodeInfo(clientID *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmEndpointInfo neofsCli.PrmEndpointInfo
	resEndpointInfo, err := cli.client.EndpointInfo(ctx, prmEndpointInfo)
	cli.mu.RUnlock()
	if err != nil {
		return errorResponsePointer("could not get endpoint info")
	}
	status := resEndpointInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	nodeInfo := resEndpointInfo.NodeInfo()
	if nodeInfo == nil {
		return errorResponsePointer("could not get node info of endpoint")
	}
	marshal, _ := nodeInfo.Marshal()
	return newResponsePointer(reflect.TypeOf(nodeInfo), marshal)
}

//export GetEndpointLatestVersion
func GetEndpointLatestVersion(clientID *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmEndpointInfo neofsCli.PrmEndpointInfo
	resEndpointInfo, err := cli.client.EndpointInfo(ctx, prmEndpointInfo)
	cli.mu.RUnlock()
	if err != nil {
		return errorResponsePointer("could not get endpoint info")
	}
	status := resEndpointInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	latestVersion := resEndpointInfo.LatestVersion()
	if latestVersion == nil {
		return errorResponsePointer("could not get latest version of endpoint")
	}
	marshal, _ := latestVersion.Marshal()
	return newResponsePointer(reflect.TypeOf(latestVersion), marshal)
}

type Endpoint struct {
	latestVersion string
	nodeInfo      string
}

func newEndpoint(latestVersion *version.Version, nodeInfo *netmap.NodeInfo) (*Endpoint, error) {
	version, err := latestVersion.Marshal()
	if err != nil {
		return nil, fmt.Errorf("could not marshal latest version of endpoint")
	}
	info, err := nodeInfo.Marshal()
	if err != nil {
		return nil, fmt.Errorf("could not marshal node info of endpoint")
	}
	return &Endpoint{
		latestVersion: string(version),
		nodeInfo:      string(info),
	}, nil
}

func endpointResponse(endpoint Endpoint) string {
	return fmt.Sprintf("{\"latestVersion\":%s,\"nodeInfo\":%s}", endpoint.latestVersion, endpoint.nodeInfo)
}

//export GetNetworkInfo
func GetNetworkInfo(clientID *C.char) C.responsePointer {
	cli, err := getClient(clientID)
	if err != nil {
		return errorResponsePointer(err.Error())
	}
	cli.mu.RLock()
	ctx := context.Background()
	var prmNetworkInfo neofsCli.PrmNetworkInfo
	resNetworkInfo, err := cli.client.NetworkInfo(ctx, prmNetworkInfo)
	cli.mu.RUnlock()
	if err != nil {
		return errorResponsePointer("could not get endpoint info")
	}
	status := resNetworkInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return resultStatusErrorResponsePointer()
	}
	info := resNetworkInfo.Info()
	if info == nil {
		return errorResponsePointer("could not get network info of endpoint")
	}
	json, err := info.Marshal()
	if err != nil {
		return errorResponsePointer("could not marshal network info of endpoint")
	}
	return newResponsePointer(reflect.TypeOf(info), json)
}
