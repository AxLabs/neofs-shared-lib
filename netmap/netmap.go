package netmap

import (
	"context"
	"encoding/json"
	"github.com/AxLabs/neofs-api-shared-lib/client"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	"github.com/google/uuid"
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

func GetEndpoint(clientID *uuid.UUID) *response.PointerResponse {
	ctx := context.Background()
	var prmEndpointInfo neofsclient.PrmEndpointInfo

	neofsClient, err := client.GetClient(clientID)
	if err != nil {
		return response.Error(err)
	}
	resEndpointInfo, err := neofsClient.LockAndGet().EndpointInfo(ctx, prmEndpointInfo)
	neofsClient.Unlock()
	if err != nil {
		return response.Error(err)
	}

	// Todo: Return specific status instead of default unsuccessful status.
	if !apistatus.IsSuccessful(resEndpointInfo.Status()) {
		return response.StatusResponse()
	}

	bytes, err := buildEndpointResponse(resEndpointInfo)
	if err != nil {
		return response.Error(err)
	}
	return response.New(reflect.TypeOf(EndpointResponse{}), bytes)
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

func GetNetworkInfo(clientID *uuid.UUID) *response.PointerResponse {
	ctx := context.Background()
	var prmNetworkInfo neofsclient.PrmNetworkInfo
	//prmNetworkInfo.WithXHeaders()

	neofsClient, err := client.GetClient(clientID)
	if err != nil {
		return response.Error(err)
	}
	resNetworkInfo, err := neofsClient.LockAndGet().NetworkInfo(ctx, prmNetworkInfo)
	neofsClient.Unlock()
	if err != nil {
		return response.Error(err)
	}

	status := resNetworkInfo.Status()
	if !apistatus.IsSuccessful(status) {
		return response.StatusResponse()
	}
	info := resNetworkInfo.Info()

	var v2 v2netmap.NetworkInfo
	info.WriteToV2(&v2)
	bytes := v2.StableMarshal(nil)

	return response.New(reflect.TypeOf(info), bytes)
}
