package gap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gap/chaincode-go/entity"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateGAP(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {

	var input entity.TransectionGAP

	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return fmt.Errorf("Unmarshal json string")
	}

	err := ctx.GetClientIdentity().AssertAttributeValue("gap.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have gap.creator role1")
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	asset := entity.TransectionGAP{
		Id:                input.Id,
		AgriStandard:      input.AgriStandard,
		FarmOwner:         input.FarmOwner,
		FarmOwnerJuristic: input.FarmOwnerJuristic,
		FarmLocation:      input.FarmLocation,
		RegisterPlants:    input.RegisterPlants,
		Owner:             clientID,
		OrgName:           orgName,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	fmt.Println(assetJSON)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, args string) error {

	var input entity.TransectionGAP
	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return fmt.Errorf("Unmarshal json string")
	}

	asset, err := s.ReadAsset(ctx, input.Id)
	if err != nil {
		return err
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if clientID != asset.Owner {
		return fmt.Errorf("submitting client not authorized to update asset, does not own asset")
	}

	asset.Id = input.Id
	asset.AgriStandard = input.AgriStandard
	asset.FarmOwner = input.FarmOwner
	asset.FarmOwnerJuristic = input.FarmOwnerJuristic
	asset.FarmLocation = input.FarmLocation
	asset.RegisterPlants = input.RegisterPlants

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

// DeleteAsset deletes a given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if clientID != asset.Owner {
		return fmt.Errorf("submitting client not authorized to update asset, does not own asset")
	}

	return ctx.GetStub().DelState(id)
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if clientID != asset.Owner {
		return fmt.Errorf("submitting client not authorized to update asset, does not own asset")
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionGAP, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionGAP
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	log.Printf("Error creating farmer chaincode: %#c", asset)

	return &asset, nil
}

func (s *SmartContract) GetAllGAP(ctx contractapi.TransactionContextInterface) ([]*entity.TransectionGAP, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	orgName, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.TransectionGAP
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionGAP
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		// Check if the asset belongs to Org1MSP
		if asset.OrgName == orgName {
			assets = append(assets, &asset)
		}
	}

	return assets, nil
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

// FilterAsset filters assets based on a specified field and value
func (s *SmartContract) FilterGap(ctx contractapi.TransactionContextInterface, typeFilter, value string) ([]*entity.TransectionGAP, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.TransectionGAP
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionGAP
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		v := reflect.ValueOf(asset)
		field := v.FieldByName(typeFilter)
		if !field.IsValid() {
			return nil, fmt.Errorf("invalid filter type: %s", typeFilter)
		}

		if field.String() == value {
			assets = append(assets, &asset)
		}
	}

	return assets, nil
}
