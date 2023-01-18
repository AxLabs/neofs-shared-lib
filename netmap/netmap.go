package netmap

import (
	"context"
	"encoding/json"
	"github.com/AxLabs/neofs-api-shared-lib/client"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	v2netmap "github.com/nspcc-dev/neofs-api-go/v2/netmap"
	v2refs "github.com/nspcc-dev/neofs-api-go/v2/refs"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"reflect"
)

/*
----Netmap----
+NetworkInfo
+EndpointInfo
*/

func GetEndpoint(neofsClient *client.NeoFSClient) *response.PointerResponse {
	ctx := context.Background()
	var prmEndpointInfo neofsclient.PrmEndpointInfo

	client := neofsClient.LockAndGet()
	resEndpointInfo, err := client.EndpointInfo(ctx, prmEndpointInfo)
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

func GetNetworkInfo(neofsClient *client.NeoFSClient) *response.PointerResponse {
	ctx := context.Background()
	var prmNetworkInfo neofsclient.PrmNetworkInfo
	//prmNetworkInfo.WithXHeaders()

	client := neofsClient.LockAndGet()
	resNetworkInfo, err := client.NetworkInfo(ctx, prmNetworkInfo)
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
