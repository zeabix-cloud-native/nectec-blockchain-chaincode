package regulator

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/regulator/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/regulator/chaincode-go/entity"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateRegulator(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityRegulator := entity.TransectionRegulator{}
	inputInterface, err := issuer.Unmarshal(args, entityRegulator)
	if err != nil {
		return err
	}
	input := inputInterface.(*entity.TransectionRegulator)

	// err := ctx.GetClientIdentity().AssertAttributeValue("regulator.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	existRegulator, err := issuer.AssetExists(ctx, input.Id)
	issuer.HandleError(err)
	if existRegulator {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	CreatedR := issuer.GetTimeNow()

	asset := entity.TransectionRegulator{
		Id:        input.Id,
		CertId:    input.CertId,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: CreatedR,
		CreatedAt: CreatedR,
	}
	assetJSON, err := json.Marshal(asset)
	issuer.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityRegulator := entity.TransectionRegulator{}
	inputInterface, err := issuer.Unmarshal(args, entityRegulator)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionRegulator)

	asset, err := s.ReadAsset(ctx, input.Id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != asset.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	UpdatedR := issuer.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedR

	assetJSON, err := json.Marshal(asset)
	issuer.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	assetRegulator, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientIDRegulator, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientIDRegulator != assetRegulator.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetR, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientIDR, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientIDR != assetR.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	assetR.Owner = newOwner
	assetJSON, err := json.Marshal(assetR)
	issuer.HandleError(err)
	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionRegulator, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionRegulator
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) GetAllRegulator(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	var filterRegulator = map[string]interface{}{}

	entityGetAll := entity.FilterGetAll{}
	interfaceRegulator, err := issuer.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := interfaceRegulator.(*entity.FilterGetAll)

	queryStringRegulator, err := issuer.BuildQueryString(filterRegulator)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryStringRegulator)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, issuer.ReturnError(issuer.SKIPOVER)
	}

	arrRegulator, err := core.FetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrRegulator, func(i, j int) bool {
		return arrRegulator[i].UpdatedAt.Before(arrRegulator[j].UpdatedAt)
	})

	if len(arrRegulator) == 0 {
		arrRegulator = []*entity.TransectionReponse{}
	}

	return &entity.GetAllReponse{
		Data:  "All Regulator",
		Obj:   arrRegulator,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterRegulator(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionRegulator, error) {
	resultsIteratorR, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIteratorR.Close()

	var assetRegulator []*entity.TransectionRegulator
	for resultsIteratorR.HasNext() {
		queryResponse, err := resultsIteratorR.Next()
		if err != nil {
			return nil, err
		}

		var dataR entity.TransectionRegulator
		err = json.Unmarshal(queryResponse.Value, &dataR)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetRegulator = append(assetRegulator, &dataR)
		}
	}

	sort.Slice(assetRegulator, func(i, j int) bool {
		return assetRegulator[i].UpdatedAt.After(assetRegulator[j].UpdatedAt)
	})

	return assetRegulator, nil
}
