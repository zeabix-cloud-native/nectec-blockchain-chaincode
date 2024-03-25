package farmer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/entity"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {

	var input entity.Transection

	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have farmer.creator role")
	}

	err := ctx.GetClientIdentity().AssertAttributeValue("farmer.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have farmer.creator role")
	}

	exists, err := s.AssetExists(ctx, input.Id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	// Get ID of submitting client identity
	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	asset := entity.Transection{
		Prefix:              input.Prefix,
		FirstName:           input.FirstName,
		LastName:            input.LastName,
		NationalID:          input.NationalID,
		AddressRegistration: input.AddressRegistration,
		Address:             input.Address,
		VillageName:         input.VillageName,
		VillageNo:           input.VillageNo,
		Road:                input.Road,
		Alley:               input.Alley,
		Subdistrict:         input.Subdistrict,
		District:            input.District,
		Province:            input.Province,
		ZipCode:             input.ZipCode,
		Phone:               input.Phone,
		MobilePhone:         input.MobilePhone,
		Email:               input.Email,
		Owner:               clientID,
		OrgName:             orgName,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string,
	firstName string,
	lastName string) error {

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

	asset.FirstName = firstName
	asset.LastName = lastName

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
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

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionReponse, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionReponse
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	log.Printf("Error creating farmer chaincode: %#c", asset)

	return &asset, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*entity.Transection, error) {

	// Query all assets from the world state
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	orgName, err := ctx.GetClientIdentity().GetMSPID()

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.Transection
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.Transection
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

// AssetExists returns true when asset with given ID exists in world state
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
func (s *SmartContract) FilterAsset(ctx contractapi.TransactionContextInterface, typeFilter, value string) ([]*entity.Transection, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.Transection
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.Transection
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		// Use reflection to get the value of the field based on typeFilter
		v := reflect.ValueOf(asset)
		field := v.FieldByName(typeFilter)
		if !field.IsValid() {
			return nil, fmt.Errorf("invalid filter type: %s", typeFilter)
		}

		// Convert field value to string and compare with the provided value
		if field.String() == value {
			assets = append(assets, &asset)
		}
	}

	return assets, nil
}
