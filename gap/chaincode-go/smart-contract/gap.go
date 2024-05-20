package gap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gap/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gap/chaincode-go/entity"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateGAP(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {

	input, err := core.UnmarshalGap(args)
	if err != nil {
		return err
	}

	// err := ctx.GetClientIdentity().AssertAttributeValue("gap.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have gap.creator role1")
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

	CreatedAt := core.GetTimeNow()

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
		Owner:       clientID,
		OrgName:     orgName,
		UpdatedAt:   CreatedAt,
		CreatedAt:   CreatedAt,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	fmt.Println(assetJSON)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, args string) error {

	input, err := core.UnmarshalGap(args)
	if err != nil {
		return err
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
		return fmt.Errorf(entity.UNAUTHORIZE)
	}

	UpdatedAt := core.GetTimeNow()

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
		return fmt.Errorf(entity.UNAUTHORIZE)
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
		return fmt.Errorf(entity.UNAUTHORIZE)
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

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
func (s *SmartContract) GetGapByCertID(ctx contractapi.TransactionContextInterface, certID string) (*entity.GetByCertIDReponse, error) {
	// Get the asset using CertID
	queryKey := fmt.Sprintf(`{"selector":{"certId":"%s"}}`, certID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryKey)
	var asset *entity.TransectionReponse
	resData := "Get gap by certID"
	if err != nil {
		return nil, fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext() {
		resData = "Not found gap by certID"

		return &entity.GetByCertIDReponse{
			Data: resData,
			Obj:  asset,
		}, nil
	}

	queryResponse, err := resultsIterator.Next()
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

	input, err := core.UnmarshalGetAll(args)
	if err != nil {
		return nil, err
	}
	filter := core.SetFilter(input)

	queryString, err := core.BuildQueryString(filter)
	if err != nil {
		return nil, err
	}

	total, err := core.CountTotalResults(ctx, queryString)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, fmt.Errorf(entity.SkipOver)
	}

	assets, err := core.FetchResultsWithPagination(ctx, input, filter)
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

func (s *SmartContract) FilterGap(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionGAP, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.TransectionGAP
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
			assets = append(assets, &asset)
		}
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].UpdatedAt.After(assets[j].UpdatedAt)
	})

	return assets, nil
}

func (s *SmartContract) CreateGapCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []entity.TransectionGAP

	errInput := json.Unmarshal([]byte(args), &inputs)
	if errInput != nil {
		return fmt.Errorf("failed to unmarshal JSON array: %v", errInput)
	}

	for _, input := range inputs {
		// err := ctx.GetClientIdentity().AssertAttributeValue("gap.creator", "true")
		// if err != nil {
		// 	return fmt.Errorf("submitting client not authorized to create asset, does not have gap.creator role1: %v", err)
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
			Owner:       clientID,
			OrgName:     orgName,
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

func (s *SmartContract) MarkUsedGap(ctx contractapi.TransactionContextInterface, args string) error {
	var input struct {
		CertId   string `json:"certId"`
		FarmerID string `json:"farmerId"`
	}
	errInput := json.Unmarshal([]byte(args), &input)
	if errInput != nil {
		return fmt.Errorf("unmarshal json string: %v", errInput)
	}

	queryKey := fmt.Sprintf(`{"selector":{"certId":"%s"}}`, input.CertId)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryKey)
	if err != nil {
		return fmt.Errorf("error querying chaincode: %v", err)
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext() {
		return fmt.Errorf("no asset found with certId: %s", input.CertId)
	}

	queryResponse, err := resultsIterator.Next()
	if err != nil {
		return fmt.Errorf("error iterating query results: %v", err)
	}

	var asset entity.TransectionGAP
	err = json.Unmarshal(queryResponse.Value, &asset)
	if err != nil {
		return fmt.Errorf("error unmarshalling asset: %v", err)
	}

	asset.FarmerID = input.FarmerID
	asset.UpdatedAt = time.Now()

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("error marshalling updated asset: %v", err)
	}

	return ctx.GetStub().PutState(asset.Id, assetJSON)
}
