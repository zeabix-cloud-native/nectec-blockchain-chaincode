/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	gap "github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gmp/chaincode-go/smart-contract"
)

func main() {
	abacSmartContract, err := contractapi.NewChaincode(&gap.SmartContract{})
	if err != nil {
		log.Panicf("Error creating gap chaincode: %v", err)
	}

	if err := abacSmartContract.Start(); err != nil {
		log.Panicf("Error starting gap chaincode: %v", err)
	}
}