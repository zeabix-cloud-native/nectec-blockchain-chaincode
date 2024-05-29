package packing

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packing/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packing/chaincode-go/entity"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreatePacking(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityPacking := entity.TransectionPacking{}
	inputInterface, err := issuer.Unmarshal(args, entityPacking)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionPacking)

	// err := ctx.GetClientIdentity().AssertAttributeValue("packing.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have packing.creator role")
	}

	existsPacking, err := issuer.AssetExists(ctx, input.Id)
	issuer.HandleError(err)
	if existsPacking {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientIDPacking, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	TimePacking := issuer.GetTimeNow()
	fmt.Printf("CreatedAt %v", TimePacking)

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
		Owner:          clientIDPacking,
		OrgName:        orgName,
		UpdatedAt:      TimePacking,
		CreatedAt:      TimePacking,
	}
	assetJSON, err := json.Marshal(asset)
	issuer.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityPacking := entity.TransectionPacking{}
	inputInterface, err := issuer.Unmarshal(args, entityPacking)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionPacking)

	asset, err := s.ReadAsset(ctx, input.Id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)
	if clientID != asset.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	UpdatedPacking := issuer.GetTimeNow()

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
	asset.UpdatedAt = UpdatedPacking

	assetJSON, errPacking := json.Marshal(asset)
	issuer.HandleError(errPacking)

	ctx.GetStub().SetEvent("UpdateAsset", assetJSON)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	assetPacking, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	clientIDPacking, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientIDPacking != assetPacking.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetPacking, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != assetPacking.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	assetPacking.Owner = newOwner
	assetJSON, err := json.Marshal(assetPacking)
	issuer.HandleError(err)
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

	return &asset, nil
}
func (s *SmartContract) GetAllPacking(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	entityGetAllPacking := entity.FilterGetAll{}
	interfacePacking, err := issuer.Unmarshal(args, entityGetAllPacking)
	if err != nil {
		return nil, err
	}
	inputPacking := interfacePacking.(*entity.FilterGetAll)
	filterPacking := core.SetFilter(inputPacking)

	queryStringPacking, err := issuer.BuildQueryString(filterPacking)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryStringPacking)
	if err != nil {
		return nil, err
	}

	if inputPacking.Skip > total {
		return nil, fmt.Errorf(issuer.SKIPOVER)
	}

	arrPacking, err := core.FetchResultsWithPagination(ctx, inputPacking, filterPacking)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrPacking, func(i, j int) bool {
		return arrPacking[i].UpdatedAt.Before(arrPacking[j].UpdatedAt)
	})

	if len(arrPacking) == 0 {
		arrPacking = []*entity.TransectionReponse{}
	}

	return &entity.GetAllReponse{
		Data:  "All Packing",
		Obj:   arrPacking,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterPacking(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionPacking, error) {
	resultsIteratorP, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIteratorP.Close()

	var assetPacking []*entity.TransectionPacking
	for resultsIteratorP.HasNext() {
		queryResponse, err := resultsIteratorP.Next()
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
			assetPacking = append(assetPacking, &asset)
		}
	}

	sort.Slice(assetPacking, func(i, j int) bool {
		return assetPacking[i].UpdatedAt.After(assetPacking[j].UpdatedAt)
	})

	return assetPacking, nil
}
