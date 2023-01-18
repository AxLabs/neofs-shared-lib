package container

import "C"
import (
	"context"
	"encoding/json"
	"github.com/AxLabs/neofs-api-shared-lib/client"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	v2container "github.com/nspcc-dev/neofs-api-go/v2/container"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	cid "github.com/nspcc-dev/neofs-sdk-go/container/id"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	"github.com/nspcc-dev/neofs-sdk-go/user"
	"reflect"
)

/*
----Container----
+Put
+Get
+Delete
+List
-SetExtendedACL
-GetExtendedACL
-AnnounceUsedSpace
*/

func PutContainer(neofsClient *client.NeoFSClient, cnr *container.Container) *response.StringResponse {
	ctx := context.Background()

	var prmContainerPut neofsclient.PrmContainerPut
	prmContainerPut.SetContainer(*cnr)

	client := neofsClient.LockAndGet()
	resContainerPut, err := client.ContainerPut(ctx, prmContainerPut)
	neofsClient.Unlock()
	if err != nil {
		return response.StringError(err)
	}

	if !apistatus.IsSuccessful(resContainerPut.Status()) {
		return response.StringStatusResponse()
	}

	containerID := resContainerPut.ID()
	return response.NewString(reflect.TypeOf(containerID), containerID.String())
}

func GetContainer(neofsClient *client.NeoFSClient, containerID *cid.ID) *response.PointerResponse {
	ctx := context.Background()

	var prmContainerGet neofsclient.PrmContainerGet
	prmContainerGet.SetContainer(*containerID)
	//prmContainerGet.WithXHeaders()

	client := neofsClient.LockAndGet()
	resContainerGet, err := client.ContainerGet(ctx, prmContainerGet)
	neofsClient.Unlock()

	if err != nil {
		return response.Error(err)
	}
	if !apistatus.IsSuccessful(resContainerGet.Status()) {
		return response.StatusResponse()
	}

	cnr := resContainerGet.Container()
	var v2 v2container.Container
	cnr.WriteToV2(&v2)
	if err != nil {
		return response.Error(err)
	}
	bytes := v2.StableMarshal(nil)
	return response.New(reflect.TypeOf(v2), bytes)
}

func DeleteContainer(neofsClient *client.NeoFSClient, containerID *cid.ID) *response.PointerResponse {
	return deleteContainer(neofsClient, containerID, nil)
}

func deleteContainer(neofsClient *client.NeoFSClient, containerID *cid.ID, sessionToken *session.Container) *response.PointerResponse {
	ctx := context.Background()

	var prmContainerDelete neofsclient.PrmContainerDelete
	prmContainerDelete.SetContainer(*containerID)
	if sessionToken != nil {
		prmContainerDelete.WithinSession(*sessionToken)
	}
	//prmContainerDelete.WithXHeaders()

	client := neofsClient.LockAndGet()
	resContainerDelete, err := client.ContainerDelete(ctx, prmContainerDelete)
	neofsClient.Unlock()
	if err != nil {
		return response.Error(err)
	}

	if !apistatus.IsSuccessful(resContainerDelete.Status()) {
		return response.StatusResponse()
	}
	return response.NewBoolean(true)
}

func ListContainer(neofsClient *client.NeoFSClient, userID *user.ID) *response.PointerResponse {
	ctx := context.Background()

	var prmContainerList neofsclient.PrmContainerList
	prmContainerList.SetAccount(*userID)
	//prmContainerList.WithXHeaders()

	client := neofsClient.LockAndGet()
	resContainerList, err := client.ContainerList(ctx, prmContainerList)
	neofsClient.Unlock()
	if err != nil {
		return response.Error(err)
	}

	if !apistatus.IsSuccessful(resContainerList.Status()) {
		return response.StatusResponse()
	}

	bytes, err := buildContainerListResponse(resContainerList)
	if err != nil {
		return response.Error(err)
	}
	return response.New(reflect.TypeOf(ListResponse{}), bytes)
}

func buildContainerListResponse(resContainerList *neofsclient.ResContainerList) ([]byte, error) {
	ids := make([]string, len(resContainerList.Containers()))
	for i, ctr := range resContainerList.Containers() {
		ids[i] = ctr.EncodeToString()
	}
	listResponse := ListResponse{Containers: ids}
	bytes, err := json.Marshal(listResponse)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type ListResponse struct {
	Containers []string `json:"containers"`
}

//endregion container
