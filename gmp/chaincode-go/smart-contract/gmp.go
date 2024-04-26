package gmp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gmp/chaincode-go/entity"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateGMP(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {

	var input entity.TransectionGMP

	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return fmt.Errorf("Unmarshal json string")
	}

	err := ctx.GetClientIdentity().AssertAttributeValue("gmp.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have gmp.creator role1")
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

	formattedTime := time.Now().Format("2006-01-02T15:04:05Z")
	CreatedAt, _ := time.Parse("2006-01-02T15:04:05Z", formattedTime)

	asset := entity.TransectionGMP{
		Id:                         input.Id,
		PackingHouseRegisterNumber: input.PackingHouseRegisterNumber,
		Address:                    input.Address,
		Owner:                      clientID,
		OrgName:                    orgName,
		UpdatedAt:                  CreatedAt,
		CreatedAt:                  CreatedAt,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	fmt.Println(assetJSON)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, args string) error {

	var input entity.TransectionGMP
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

	formattedTime := time.Now().Format("2006-01-02T15:04:05Z")
	UpdatedAt, _ := time.Parse("2006-01-02T15:04:05Z", formattedTime)

	asset.Id = input.Id
	asset.PackingHouseRegisterNumber = input.PackingHouseRegisterNumber
	asset.Address = input.Address
	asset.UpdatedAt = UpdatedAt

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

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionGMP, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionGMP
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	log.Printf("Error creating farmer chaincode: %#c", asset)

	return &asset, nil
}

func (s *SmartContract) GetAllGMP(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	var input entity.FilterGetAll
	var filter = map[string]interface{}{}

	errInput := json.Unmarshal([]byte(args), &input)
	if errInput != nil {
		return nil, fmt.Errorf("Unmarshal json string")
	}

	if input.PackingHouseRegisterNumber != nil {
		filter["packingHouseRegisterNumber"] = input.PackingHouseRegisterNumber
	}

	if input.Address != nil {
		filter["address"] = input.Address
	}

	limit := input.Limit
	skip := input.Skip

	selector := map[string]interface{}{
		"selector": filter,
	}

	queryString, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryString))
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	total := 0

	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		total++
	}
	if skip > total {
		return nil, fmt.Errorf("Skip Over Total Data")
	}
	// Check if skip and limit are provided in the args
	if skip != 0 || limit != 0 {
		selector["skip"] = skip
		selector["limit"] = limit
	}

	queryStringPagination, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryResults, _, err := ctx.GetStub().GetQueryResultWithPagination(string(queryStringPagination), int32(limit), "")
	if err != nil {
		return nil, err
	}
	defer queryResults.Close()

	var assets []*entity.TransectionReponse

	for queryResults.HasNext() {
		queryResponse, err := queryResults.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionReponse
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].UpdatedAt.Before(assets[j].UpdatedAt)
	})

	if len(assets) == 0 {
		assets = []*entity.TransectionReponse{}
	}

	return &entity.GetAllReponse{
		Data:  "All Gmp",
		Obj:   assets,
		Total: total,
	}, nil
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

func (s *SmartContract) FilterGmp(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionGMP, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.TransectionGMP
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionGMP
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assets = append(assets, &asset)
		}
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].UpdatedAt.After(assets[j].UpdatedAt)
	})

	return assets, nil
}

func (s *SmartContract) CreateGmpCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []entity.TransectionGMP

	errInput := json.Unmarshal([]byte(args), &inputs)
	if errInput != nil {
		return fmt.Errorf("failed to unmarshal JSON array: %v", errInput)
	}

	for _, input := range inputs {
		// err := ctx.GetClientIdentity().AssertAttributeValue("gmp.creator", "true")
		// if err != nil {
		// 	return fmt.Errorf("submitting client not authorized to create asset, does not have gmp.creator role1: %v", err)
		// }

		orgName, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		exists, err := s.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if exists {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		clientID, err := s.GetSubmittingClientIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		asset := entity.TransectionGMP{
			Id:                         input.Id,
			PackingHouseRegisterNumber: input.PackingHouseRegisterNumber,
			Address:                    input.Address,
			Owner:                      clientID,
			OrgName:                    orgName,
		}
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return fmt.Errorf("failed to marshal asset JSON: %v", err)
		}

		err = ctx.GetStub().PutState(input.Id, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
		}

		fmt.Printf("Asset %s created successfully\n", input.Id)
	}

	return nil
}
