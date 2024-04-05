/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	packing "github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packing/chaincode-go/smart-contract"
)

func main() {
	abacSmartContract, err := contractapi.NewChaincode(&packing.SmartContract{})
	if err != nil {
		log.Panicf("Error creating packing chaincode: %v", err)
	}

	if err := abacSmartContract.Start(); err != nil {
		log.Panicf("Error starting packing chaincode: %v", err)
	}
}
