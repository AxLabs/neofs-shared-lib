package accounting

import (
	"context"
	"github.com/AxLabs/neofs-api-shared-lib/client"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	"github.com/google/uuid"
	v2accounting "github.com/nspcc-dev/neofs-api-go/v2/accounting"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	apistatus "github.com/nspcc-dev/neofs-sdk-go/client/status"
	"github.com/nspcc-dev/neofs-sdk-go/user"
	"reflect"
)

/*
----Accounting----
Balance
*/

func GetBalance(clientID *uuid.UUID, id *user.ID) *response.PointerResponse {
	ctx := context.Background()

	var prmBalanceGet neofsclient.PrmBalanceGet
	//id, err := main.UserIDFromPublicKey(publicKey)
	prmBalanceGet.SetAccount(*id)

	neofsClient, err := client.GetClient(clientID)
	if err != nil {
		return response.ClientError()
	}
	resBalanceGet, err := neofsClient.LockAndGet().BalanceGet(ctx, prmBalanceGet)
	neofsClient.Unlock()
	if err != nil {
		return response.Error(err)
	}

	resStatus := resBalanceGet.Status()
	if !apistatus.IsSuccessful(resStatus) {
		return response.StatusResponse()
	}

	amount := resBalanceGet.Amount()
	if amount == nil {
		return response.Error(err)
	}

	var v2 v2accounting.Decimal
	amount.WriteToV2(&v2)
	bytes := v2.StableMarshal(nil)
	return response.New(reflect.TypeOf(v2), bytes)
}
