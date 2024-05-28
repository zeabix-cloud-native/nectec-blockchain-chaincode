package packer

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packer/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packer/chaincode-go/entity"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreatePacker(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityPacker := entity.TransectionPacker{}
	inputInterface, err := issuer.Unmarshal(args, entityPacker)
	if err != nil {
		return err
	}
	input := inputInterface.(*entity.TransectionPacker)

	// err := ctx.GetClientIdentity().AssertAttributeValue("packer.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have packer.creator role")
	}

	existPacker, err := issuer.AssetExists(ctx, input.Id)
	issuer.HandleError(err)
	if existPacker {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	TimePacker := issuer.GetTimeNow()

	asset := entity.TransectionPacker{
		Id:        input.Id,
		CertId:    input.CertId,
		UserId:    input.UserId,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: TimePacker,
		CreatedAt: TimePacker,
	}
	assetJSON, err := json.Marshal(asset)
	issuer.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityPacker := entity.TransectionPacker{}
	inputInterface, err := issuer.Unmarshal(args, entityPacker)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionPacker)

	asset, err := s.ReadAsset(ctx, input.Id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != asset.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	UpdatedPacker := issuer.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UserId = input.UserId
	asset.UpdatedAt = UpdatedPacker

	assetJSON, errP := json.Marshal(asset)
	issuer.HandleError(errP)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	assetPacker, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientIDPacker, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientIDPacker != assetPacker.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetP, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != assetP.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	assetP.Owner = newOwner
	assetJSON, err := json.Marshal(assetP)
	issuer.HandleError(err)
	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionPacker, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionPacker
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) GetPackerById(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionReponse, error) {
	queryPacker := fmt.Sprintf(`{"selector":{"id":"%s"}}`, id)

	resultsPacker, err := ctx.GetStub().GetQueryResult(queryPacker)
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsPacker.Close()

	if !resultsPacker.HasNext() {
		return nil, fmt.Errorf("the asset with id %s does not exist", id)
	}

	queryResponse, err := resultsPacker.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	var asset entity.TransectionReponse
	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &asset, nil
}

func (s *SmartContract) GetAllPacker(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	var filterPacker = map[string]interface{}{}

	entityGetAll := entity.FilterGetAll{}
	interfacePacker, err := issuer.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := interfacePacker.(*entity.FilterGetAll)

	queryStringPacker, err := issuer.BuildQueryString(filterPacker)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryStringPacker)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, issuer.ReturnError(issuer.SKIPOVER)
	}

	arrPacker, err := core.FetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrPacker, func(i, j int) bool {
		return arrPacker[i].UpdatedAt.Before(arrPacker[j].UpdatedAt)
	})

	if len(arrPacker) == 0 {
		arrPacker = []*entity.TransectionReponse{}
	}

	return &entity.GetAllReponse{
		Data:  "All Packer",
		Obj:   arrPacker,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterPacker(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionPacker, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetPacker []*entity.TransectionPacker
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionPacker
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetPacker = append(assetPacker, &asset)
		}
	}

	sort.Slice(assetPacker, func(i, j int) bool {
		return assetPacker[i].UpdatedAt.After(assetPacker[j].UpdatedAt)
	})

	return assetPacker, nil
}
