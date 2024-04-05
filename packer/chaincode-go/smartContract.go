/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	packer "github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packer/chaincode-go/smart-contract"
)

func main() {
	abacSmartContract, err := contractapi.NewChaincode(&packer.SmartContract{})
	if err != nil {
		log.Panicf("Error creating nstda staff chaincode: %v", err)
	}

	if err := abacSmartContract.Start(); err != nil {
		log.Panicf("Error starting nstda staff chaincode: %v", err)
	}
}
