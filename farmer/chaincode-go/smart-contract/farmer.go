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
		return fmt.Errorf("Unmarshal json string")
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

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {

	var input entity.Transection
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
	asset.Prefix = input.Prefix
	asset.FirstName = input.FirstName
	asset.LastName = input.LastName
	asset.NationalID = input.NationalID
	asset.AddressRegistration = input.AddressRegistration
	asset.Address = input.Address
	asset.VillageName = input.VillageName
	asset.VillageNo = input.VillageNo
	asset.Road = input.Road
	asset.Alley = input.Alley
	asset.Subdistrict = input.Subdistrict
	asset.District = input.District
	asset.Province = input.Province
	asset.ZipCode = input.ZipCode
	asset.Phone = input.Phone
	asset.MobilePhone = input.MobilePhone
	asset.Email = input.Email

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

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.Transection, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.Transection
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	log.Printf("Error creating farmer chaincode: %#c", asset)

	return &asset, nil
}

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

func (s *SmartContract) FilterFarmer(ctx contractapi.TransactionContextInterface, typeFilter, value string) ([]*entity.Transection, error) {
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
