package client

import "C"
import (
	"crypto/ecdsa"
	"fmt"
	"github.com/AxLabs/neofs-api-shared-lib/response"
	"github.com/google/uuid"
	neofsclient "github.com/nspcc-dev/neofs-sdk-go/client"
	"reflect"
	"sync"
)

var neofsClientMap *NeoFSClientMap

type NeoFSClient struct {
	mu     sync.RWMutex
	client *neofsclient.Client
}

type NeoFSClientMap struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]*NeoFSClient
}

func initClientMap(id uuid.UUID, newClient *neofsclient.Client) {
	neofsClientMap = &NeoFSClientMap{sync.RWMutex{}, map[uuid.UUID]*NeoFSClient{id: {sync.RWMutex{}, newClient}}}
}

func (clients *NeoFSClientMap) put(id uuid.UUID, newClient *neofsclient.Client) {
	clients.mu.Lock()
	clients.clients[id] = &NeoFSClient{
		mu:     sync.RWMutex{},
		client: newClient,
	}
	clients.mu.Unlock()
}

func (clients *NeoFSClientMap) delete(id uuid.UUID) bool {
	clients.mu.Lock()
	delete(clients.clients, id)
	clients.mu.Unlock()
	return true
}

func (c *NeoFSClient) LockAndGet() *neofsclient.Client {
	c.mu.Lock()
	return c.client
}

func (c *NeoFSClient) Unlock() {
	c.mu.Unlock()
}

func GetClient(clientID *uuid.UUID) (*NeoFSClient, error) {
	if neofsClientMap == nil {
		return nil, fmt.Errorf("no clients present")
	}
	neofsClientMap.mu.Lock()
	cli := neofsClientMap.clients[*clientID]
	if cli == nil {
		return nil, fmt.Errorf("no client present with id %v", clientID)
	}
	neofsClientMap.mu.Unlock()
	return cli, nil
}

func CreateClient(privateKey *ecdsa.PrivateKey, neofsEndpoint string) *response.PointerResponse {
	// Initialize client
	newClient := neofsclient.Client{}
	var prmInit neofsclient.PrmInit
	prmInit.SetDefaultPrivateKey(*privateKey)
	prmInit.ResolveNeoFSFailures()
	//prmInit.SetResponseInfoCallback()
	newClient.Init(prmInit)

	// Set dial configuration in client
	var prmDial neofsclient.PrmDial
	prmDial.SetServerURI(neofsEndpoint)
	//prmDial.SetTLSConfig() // default means insecure connection
	//prmDial.SetTimeout() // 5 seconds by default
	err := newClient.Dial(prmDial)
	if err != nil {
		return response.Error(err)
	}

	u, err := uuid.NewUUID()
	if err != nil {
		return response.Error(err)
	}

	if neofsClientMap == nil {
		initClientMap(u, &newClient)
	} else {
		neofsClientMap.put(u, &newClient)
	}
	return response.New(reflect.TypeOf(u), []byte(u.String()))
}

////export DeleteClient
//func DeleteClient(clientID *C.char) C.PointerResponse {
//	cliID, err := uuid.Parse(C.GoString(clientID))
//	if err != nil {
//		return PointerResponseError("could not parse provided client id")
//	}
//	deleted := neofsClients.delete(cliID)
//	if !deleted {
//		return PointerResponseError("could not delete client")
//	}
//	boolean := []byte{1}
//	return PointerResponse(reflect.TypeOf(boolean), boolean)
//}
