package gmp

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gmp/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gmp/chaincode-go/entity"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateGMP(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityGmp := entity.TransectionGMP{}
	inputInterface, err := issuer.Unmarshal(args, entityGmp)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionGMP)

	// err := ctx.GetClientIdentity().AssertAttributeValue("gmp.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have gmp.creator role1")
	}

	existsGmp, err := issuer.AssetExists(ctx, input.Id)
	issuer.HandleError(err)
	if existsGmp {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	TimeGmp := issuer.GetTimeNow()

	asset := entity.TransectionGMP{
		Id:                         input.Id,
		PackerId: 									input.PackerId,		
		PackingHouseRegisterNumber: input.PackingHouseRegisterNumber,
		Address:                    input.Address,
		PackingHouseName:           input.PackingHouseName,
		UpdatedDate:                input.UpdatedDate,
		Source:                     input.Source,
		Owner:                      clientID,
		OrgName:                    orgName,
		UpdatedAt:                  TimeGmp,
		CreatedAt:                  TimeGmp,
	}
	assetJSON, err := json.Marshal(asset)
	issuer.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, args string) error {

	entityGmp := entity.TransectionGMP{}
	inputInterface, err := issuer.Unmarshal(args, entityGmp)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionGMP)

	asset, err := s.ReadAsset(ctx, input.Id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)
	if clientID != asset.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	UpdatedGmp := issuer.GetTimeNow()

	asset.Id = input.Id
	asset.PackerId = input.PackerId
	asset.PackingHouseRegisterNumber = input.PackingHouseRegisterNumber
	asset.Address = input.Address
	asset.PackingHouseName = input.PackingHouseName
	asset.UpdatedDate = input.UpdatedDate
	asset.Source = input.Source
	asset.UpdatedAt = UpdatedGmp

	assetJSON, errG := json.Marshal(asset)
	issuer.HandleError(errG)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	assetGmp, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientIDGmp, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientIDGmp != assetGmp.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetG, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != assetG.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	assetG.Owner = newOwner
	assetJSON, err := json.Marshal(assetG)
	issuer.HandleError(err)
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

	return &asset, nil
}

func (s *SmartContract) GetAllGMP(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	entityGetAllGmp := entity.FilterGetAll{}
	interfaceGmp, err := issuer.Unmarshal(args, entityGetAllGmp)
	if err != nil {
		return nil, err
	}
	inputGmp := interfaceGmp.(*entity.FilterGetAll)
	filterGmp := core.SetFilter(inputGmp)

	queryStringGmp, err := issuer.BuildQueryString(filterGmp)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryStringGmp)
	if err != nil {
		return nil, err
	}

	if inputGmp.Skip > total {
		return nil, fmt.Errorf(issuer.SKIPOVER)
	}

	assets, err := core.FetchResultsWithPagination(ctx, inputGmp, filterGmp)
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
		Data:  "All Gmp",
		Obj:   assets,
		Total: total,
	}, nil
}

func (s *SmartContract) GetGmpByPackingHouseNumber(ctx contractapi.TransactionContextInterface, packingHouseRegisterNumber string) (*entity.GetByRegisterNumberResponse, error) {
	// Get the asset using CertID
	queryKeyPackingHouse := fmt.Sprintf(`{"selector":{"packingHouseRegisterNumber":"%s"}}`, packingHouseRegisterNumber)

	resultsIteratorPackingHouse, err := ctx.GetStub().GetQueryResult(queryKeyPackingHouse)
	var asset *entity.TransectionReponse
	resData := "Get gmp by packingHouseRegisterNumber"
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIteratorPackingHouse.Close()

	if !resultsIteratorPackingHouse.HasNext() {
		resData = "Not found gmp by packingHouseRegisterNumber"

		return &entity.GetByRegisterNumberResponse{
			Data: resData,
			Obj:  asset,
		}, nil
	}

	queryResponse, err := resultsIteratorPackingHouse.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &entity.GetByRegisterNumberResponse{
		Data: resData,
		Obj:  asset,
	}, nil

}

func (s *SmartContract) FilterGmp(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionGMP, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetGmp []*entity.TransectionGMP
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
			assetGmp = append(assetGmp, &asset)
		}
	}

	sort.Slice(assetGmp, func(i, j int) bool {
		return assetGmp[i].UpdatedAt.After(assetGmp[j].UpdatedAt)
	})

	return assetGmp, nil
}

func (s *SmartContract) CreateGmpCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []entity.TransectionGMP

	errInputGmp := json.Unmarshal([]byte(args), &inputs)
	issuer.HandleError(errInputGmp)

	for _, input := range inputs {
		// err := ctx.GetClientIdentity().AssertAttributeValue("gmp.creator", "true")

		orgNameG, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		existGmp, err := issuer.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existGmp {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		clientIDG, err := issuer.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		assetG := entity.TransectionGMP{
			Id:                         input.Id,
			PackerId: 									input.PackerId,
			PackingHouseRegisterNumber: input.PackingHouseRegisterNumber,
			Address:                    input.Address,
			PackingHouseName:           input.PackingHouseName,
			UpdatedDate:                input.UpdatedDate,
			Source:                     input.Source,
			Owner:                      clientIDG,
			OrgName:                    orgNameG,
		}
		assetJSON, err := json.Marshal(assetG)
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
