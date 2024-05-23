package exporter

import (
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

	existExporter, err := issuer.AssetExists(ctx, input.Id)
	if err != nil {
		return err
	}
	if existExporter {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

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
	issuer.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityExporter := entity.TransectionExporter{}
	inputInterface, err := issuer.Unmarshal(args, entityExporter)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionExporter)

	asset, err := s.ReadAsset(ctx, input.Id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != asset.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	UpdatedTime := issuer.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedTime

	assetJSON, errE := json.Marshal(asset)
	issuer.HandleError(errE)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	assetE, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientIDExporter, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientIDExporter != assetE.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetE, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != assetE.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	assetE.Owner = newOwner
	assetJSON, err := json.Marshal(assetE)
	issuer.HandleError(err)
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

	var filterE = map[string]interface{}{}

	entityGetAll := entity.FilterGetAll{}
	interfaceE, err := issuer.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := interfaceE.(*entity.FilterGetAll)

	queryStringE, err := issuer.BuildQueryString(filterE)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryStringE)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, fmt.Errorf(issuer.SKIPOVER)
	}

	arrExporter, err := core.FetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrExporter, func(i, j int) bool {
		return arrExporter[i].UpdatedAt.Before(arrExporter[j].UpdatedAt)
	})

	if len(arrExporter) == 0 {
		arrExporter = []*entity.TransectionReponse{}
	}

	return &entity.GetAllReponse{
		Data:  "All Exporter",
		Obj:   arrExporter,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterExporter(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionExporter, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetExporter []*entity.TransectionExporter
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
			assetExporter = append(assetExporter, &asset)
		}
	}

	sort.Slice(assetExporter, func(i, j int) bool {
		return assetExporter[i].UpdatedAt.After(assetExporter[j].UpdatedAt)
	})

	return assetExporter, nil
}
