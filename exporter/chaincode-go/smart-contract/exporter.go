package exporter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/exporter/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/exporter/chaincode-go/entity"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateExporter(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityExporter := entity.TransectionExporter{}
	inputInterface, err := issuer.Unmarshal(args, entityExporter)
	if err != nil {
		return err
	}
	input := inputInterface.(*entity.TransectionExporter)

	// err := ctx.GetClientIdentity().AssertAttributeValue("exporter.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have exporter.creator role")
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

	CreatedTime := issuer.GetTimeNow()

	asset := entity.TransectionExporter{
		Id:        input.Id,
		CertId:    input.CertId,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: CreatedTime,
		CreatedAt: CreatedTime,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityExporter := entity.TransectionExporter{}
	inputInterface, err := issuer.Unmarshal(args, entityExporter)
	if err != nil {
		return err
	}
	input := inputInterface.(*entity.TransectionExporter)

	asset, err := s.ReadAsset(ctx, input.Id)
	if err != nil {
		return err
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if clientID != asset.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	UpdatedTime := issuer.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedTime

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
		return fmt.Errorf(issuer.UNAUTHORIZE)
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
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionExporter, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionExporter
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) GetAllExporter(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	var filter = map[string]interface{}{}

	entityGetAll := entity.FilterGetAll{}
	inputInterface, err := issuer.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := inputInterface.(*entity.FilterGetAll)

	queryString, err := issuer.BuildQueryString(filter)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryString)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, fmt.Errorf(issuer.SKIPOVER)
	}

	assets, err := core.FetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].UpdatedAt.Before(assets[j].UpdatedAt)
	})

	if len(assets) == 0 {
		assets = []*entity.TransectionReponse{}
	}

	return &entity.GetAllReponse{
		Data:  "All Exporter",
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
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func (s *SmartContract) FilterExporter(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionExporter, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.TransectionExporter
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionExporter
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
