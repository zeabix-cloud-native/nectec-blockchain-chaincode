package packing

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packing/chaincode-go/entity"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreatePacking(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {

	var input entity.TransectionPacking

	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return fmt.Errorf("unmarshal json string")
	}

	err := ctx.GetClientIdentity().AssertAttributeValue("packing.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have packing.creator role")
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

	asset := entity.TransectionPacking{
		Id:             input.Id,
		OrderID:        input.OrderID,
		FarmerID:       input.FarmerID,
		ForecastWeight: input.ForecastWeight,
		ActualWeight:   input.ActualWeight,
		SavedTime:      input.SavedTime,
		ApprovedDate:   input.ApprovedDate,
		ApprovedType:   input.ApprovedType,
		FinalWeight:    input.FinalWeight,
		Remark:         input.Remark,
		PackerId:       input.PackerId,
		Gmp:            input.Gmp,
		Gap:            input.Gap,
		ProcessStatus:  input.ProcessStatus,
		Owner:          clientID,
		OrgName:        orgName,
		UpdatedAt:      CreatedAt,
		CreatedAt:      CreatedAt,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {

	var input entity.TransectionPacking
	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return fmt.Errorf("unmarshal json string")
	}

	asset, err := s.ReadAsset(ctx, input.Id)
	if err != nil {
		return err
	}

	if asset.FarmerID != input.FarmerID || asset.PackerId != input.PackerId {
		return fmt.Errorf("FarmerID or PackerId not match")
	}

	formattedTime := time.Now().Format("2006-01-02T15:04:05Z")
	UpdatedAt, _ := time.Parse("2006-01-02T15:04:05Z", formattedTime)

	asset.Id = input.Id
	asset.OrderID = input.OrderID
	asset.FarmerID = input.FarmerID // not update
	asset.ForecastWeight = input.ForecastWeight
	asset.ActualWeight = input.ActualWeight
	asset.SavedTime = input.SavedTime
	asset.ApprovedDate = input.ApprovedDate
	asset.ApprovedType = input.ApprovedType
	asset.FinalWeight = input.FinalWeight
	asset.Remark = input.Remark
	asset.PackerId = input.PackerId // not update
	asset.Gmp = input.Gmp
	asset.Gap = input.Gap
	asset.ProcessStatus = input.ProcessStatus
	asset.UpdatedAt = UpdatedAt

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("UpdateAsset", assetJSON)

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

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionPacking, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionPacking
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	log.Printf("Error creating farmer chaincode: %#c", asset)

	return &asset, nil
}
func (s *SmartContract) GetAllPacking(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	var input entity.FilterGetAll
	var filter = map[string]interface{}{}

	errInput := json.Unmarshal([]byte(args), &input)
	if errInput != nil {
		return nil, fmt.Errorf("unmarshal json string")
	}

	if input.PackerId != nil {
		filter["packerId"] = *input.PackerId
	}

	if input.Gap != nil {
		filter["gap"] = *input.Gap
	}

	if input.StartDate != nil && input.EndDate != nil {
		filter["createdAt"] = map[string]interface{}{
			"$gte": *input.StartDate,
			"$lte": *input.EndDate,
		}
	}

	if input.ForecastWeightFrom != nil && input.ForecastWeightTo != nil {
		filter["forecastWeight"] = map[string]interface{}{
			"$gte": *input.ForecastWeightFrom,
			"$lte": *input.ForecastWeightTo,
		}
	}

	if input.ProcessStatus != nil {
		filter["processStatus"] = *input.ProcessStatus
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
		Data:  "All Packing",
		Obj:   assets,
		Total: total,
	}, nil
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
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func (s *SmartContract) FilterPacking(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionPacking, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.TransectionPacking
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionPacking
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
