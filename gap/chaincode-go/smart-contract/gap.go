package gap

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gap/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gap/chaincode-go/entity"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateGAP(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityGap := entity.TransectionGAP{}
	inputInterface, err := issuer.Unmarshal(args, entityGap)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionGAP)

	// err := ctx.GetClientIdentity().AssertAttributeValue("gap.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have gap.creator role1")
	}

	existsGap, err := issuer.AssetExists(ctx, input.Id)
	issuer.HandleError(err)
	if existsGap {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientIDGap, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	TimeGap := issuer.GetTimeNow()

	asset := entity.TransectionGAP{
		Id:          input.Id,
		CertID:      input.CertID,
		AreaCode:    input.AreaCode,
		AreaRai:     input.AreaRai,
		AreaStatus:  input.AreaStatus,
		OldAreaCode: input.OldAreaCode,
		IssueDate:   input.IssueDate,
		ExpireDate:  input.ExpireDate,
		District:    input.District,
		Province:    input.Province,
		UpdatedDate: input.UpdatedDate,
		Source:      input.Source,
		FarmerID:    input.FarmerID,
		Owner:       clientIDGap,
		OrgName:     orgName,
		UpdatedAt:   TimeGap,
		CreatedAt:   TimeGap,
	}
	assetJSON, err := json.Marshal(asset)
	issuer.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, args string) error {

	entityGap := entity.TransectionGAP{}
	inputInterface, err := issuer.Unmarshal(args, entityGap)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionGAP)

	asset, err := s.ReadAsset(ctx, input.Id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)
	if clientID != asset.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	UpdatedGap := issuer.GetTimeNow()

	asset.Id = input.Id
	asset.CertID = input.CertID
	asset.AreaCode = input.AreaCode
	asset.AreaRai = input.AreaRai
	asset.AreaStatus = input.AreaStatus
	asset.OldAreaCode = input.OldAreaCode
	asset.IssueDate = input.IssueDate
	asset.ExpireDate = input.ExpireDate
	asset.District = input.District
	asset.Province = input.Province
	asset.UpdatedDate = input.UpdatedDate
	asset.Source = input.Source
	asset.FarmerID = input.FarmerID
	asset.UpdatedAt = UpdatedGap

	assetJSON, errGap := json.Marshal(asset)
	issuer.HandleError(errGap)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	assetGap, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientIDGap, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientIDGap != assetGap.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetGap, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != assetGap.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	assetGap.Owner = newOwner
	assetJSON, err := json.Marshal(assetGap)
	issuer.HandleError(err)
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

	return &asset, nil
}

func (s *SmartContract) GetGapByFarmerID(ctx contractapi.TransactionContextInterface, farmerId string) (*entity.GetByCertIDReponse, error) {
	// Get the asset using farmerId 
	queryKeyFarmer := fmt.Sprintf(`{"selector":{"farmerId":"%s"}}`, farmerId)

	resultsIteratorFarmer, err := ctx.GetStub().GetQueryResult(queryKeyFarmer)
	var asset *entity.TransectionReponse
	resData := "Get gap by farmerId"
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIteratorFarmer.Close()

	if !resultsIteratorFarmer.HasNext() {
		resData = "Not found gap by farmerId"

		return &entity.GetByCertIDReponse{
			Data: resData,
			Obj:  asset,
		}, nil
	}

	queryResponse, err := resultsIteratorFarmer.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &entity.GetByCertIDReponse{
		Data: resData,
		Obj:  asset,
	}, nil

}

func (s *SmartContract) GetGapByCertID(ctx contractapi.TransactionContextInterface, certID string) (*entity.GetByCertIDReponse, error) {
	// Get the asset using CertID
	queryKeyGap := fmt.Sprintf(`{"selector":{"certId":"%s"}}`, certID)

	resultsIteratorGap, err := ctx.GetStub().GetQueryResult(queryKeyGap)
	var asset *entity.TransectionReponse
	resData := "Get gap by certID"
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIteratorGap.Close()

	if !resultsIteratorGap.HasNext() {
		resData = "Not found gap by certID"

		return &entity.GetByCertIDReponse{
			Data: resData,
			Obj:  asset,
		}, nil
	}

	queryResponse, err := resultsIteratorGap.Next()
	if err != nil {
		return nil, fmt.Errorf("error getting next query result: %v", err)
	}

	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling asset JSON: %v", err)
	}

	return &entity.GetByCertIDReponse{
		Data: resData,
		Obj:  asset,
	}, nil

}

func (s *SmartContract) GetAllGAP(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	entityGetAllGap := entity.FilterGetAll{}
	interfaceGap, err := issuer.Unmarshal(args, entityGetAllGap)
	if err != nil {
		return nil, err
	}
	inputGap := interfaceGap.(*entity.FilterGetAll)
	filterGap := core.SetFilter(inputGap)

	queryStringGap, err := issuer.BuildQueryString(filterGap)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryStringGap)
	if err != nil {
		return nil, err
	}

	if inputGap.Skip > total {
		return nil, fmt.Errorf(issuer.SKIPOVER)
	}

	assets, err := core.FetchResultsWithPagination(ctx, inputGap, filterGap)
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
		Data:  "All Gap",
		Obj:   assets,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterGap(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionGAP, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetGap []*entity.TransectionGAP
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

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetGap = append(assetGap, &asset)
		}
	}

	sort.Slice(assetGap, func(i, j int) bool {
		return assetGap[i].UpdatedAt.After(assetGap[j].UpdatedAt)
	})

	return assetGap, nil
}

func (s *SmartContract) CreateGapCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []entity.TransectionGAP

	errInputGap := json.Unmarshal([]byte(args), &inputs)
	issuer.HandleError(errInputGap)

	for _, input := range inputs {
		// err := ctx.GetClientIdentity().AssertAttributeValue("gap.creator", "true")

		orgNameGap, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return fmt.Errorf("failed to get submitting client's MSP ID: %v", err)
		}

		existGap, err := issuer.AssetExists(ctx, input.Id)
		if err != nil {
			return fmt.Errorf("error checking if asset exists: %v", err)
		}
		if existGap {
			return fmt.Errorf("the asset %s already exists", input.Id)
		}

		clientIDGap, err := issuer.GetIdentity(ctx)
		if err != nil {
			return fmt.Errorf("failed to get submitting client's identity: %v", err)
		}

		assetGap := entity.TransectionGAP{
			Id:          input.Id,
			CertID:      input.CertID,
			AreaCode:    input.AreaCode,
			AreaRai:     input.AreaRai,
			AreaStatus:  input.AreaStatus,
			OldAreaCode: input.OldAreaCode,
			IssueDate:   input.IssueDate,
			ExpireDate:  input.ExpireDate,
			District:    input.District,
			Province:    input.Province,
			UpdatedDate: input.UpdatedDate,
			Source:      input.Source,
			FarmerID:    input.FarmerID,
			Owner:       clientIDGap,
			OrgName:     orgNameGap,
		}
		assetJSON, err := json.Marshal(assetGap)
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
