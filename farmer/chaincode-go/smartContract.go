/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	farmer "github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/smart-contract"
)

func main() {
	abacSmartContract, err := contractapi.NewChaincode(&farmer.SmartContract{})
	if err != nil {
		log.Panicf("Error creating farmer chaincode: %v", err)
	}

	if err := abacSmartContract.Start(); err != nil {
		log.Panicf("Error starting farmer chaincode: %v", err)
	}
}
