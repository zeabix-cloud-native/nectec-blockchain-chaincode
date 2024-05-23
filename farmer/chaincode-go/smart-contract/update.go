package farmer

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/entity"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
)

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {
	entityType := entity.TransectionFarmer{}
	inputInterface, err := issuer.Unmarshal(args, entityType)
	if err != nil {
		return err
	}
	input := inputInterface.(*entity.TransectionFarmer)

	asset, err := s.ReadAsset(ctx, input.Id)
	if err != nil {
		return err
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if clientID != asset.Owner {
		return fmt.Errorf(entity.UNAUTHORIZE)
	}

	UpdatedAt := issuer.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedAt

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(input.Id, assetJSON)
}
