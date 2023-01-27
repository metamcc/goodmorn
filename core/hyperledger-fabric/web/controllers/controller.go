package controllers

import (
	"github.com/mcc/hyperledger-fabric-server/blockchain"
)

type Application struct {
	Fabric *blockchain.FabricSetup
}

type FabricData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
