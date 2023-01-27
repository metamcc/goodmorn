package main

import (
	"fmt"
	"os"

	"github.com/mcc/hyperledger-fabric-server/blockchain"
	"github.com/mcc/hyperledger-fabric-server/web"
	"github.com/mcc/hyperledger-fabric-server/web/controllers"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Channel parameters
		ChannelID:     "mycreditchain",
		ChannelConfig: "/home/vagrant/go/src/github.com/mcc/hyperledger-fabric-server/fixtures_test/channel-artifacts/mycreditchain.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "mycreditchain-service",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/mcc/hyperledger-fabric-server/chaincode/",
		OrgAdmin:        "admin.mccorg",
		OrgName:         "mccorg",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "user.mccorg",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
	}

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
	}

	// Launch the web application listening
	app := &controllers.Application{
		Fabric: &fSetup,
	}
	web.Serve(app)
}
