package nstdaStaff

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/nstda-staff/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/nstda-staff/chaincode-go/entity"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateNstdaStaff(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityNstda := entity.TransectionNstdaStaff{}
	inputInterface, err := issuer.Unmarshal(args, entityNstda)
	if err != nil {
		return err
	}
	input := inputInterface.(*entity.TransectionNstdaStaff)

	// err := ctx.GetClientIdentity().AssertAttributeValue("nstdaStaff.creator", "true")
	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have nstdaStaff.creator role")
	}

	existNstda, err := issuer.AssetExists(ctx, input.Id)
	issuer.HandleError(err)
	if existNstda {
		return fmt.Errorf("the asset %s already exists", input.Id)
	}

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	TimeNstda := issuer.GetTimeNow()

	asset := entity.TransectionNstdaStaff{
		Id:        input.Id,
		CertId:    input.CertId,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: TimeNstda,
		CreatedAt: TimeNstda,
	}
	assetJSON, err := json.Marshal(asset)
	issuer.HandleError(err)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {

	entityNstda := entity.TransectionNstdaStaff{}
	inputInterface, err := issuer.Unmarshal(args, entityNstda)
	issuer.HandleError(err)
	input := inputInterface.(*entity.TransectionNstdaStaff)

	asset, err := s.ReadAsset(ctx, input.Id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != asset.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	UpdatedNstda := issuer.GetTimeNow()
	asset.Id = input.Id
	asset.CertId = input.CertId
	asset.UpdatedAt = UpdatedNstda

	assetJSON, errN := json.Marshal(asset)
	issuer.HandleError(errN)

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	assetNstda, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientIDNstda, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientIDNstda != assetNstda.Owner {
		return fmt.Errorf(issuer.UNAUTHORIZE)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {

	assetN, err := s.ReadAsset(ctx, id)
	issuer.HandleError(err)

	clientID, err := issuer.GetIdentity(ctx)
	issuer.HandleError(err)

	if clientID != assetN.Owner {
		return issuer.ReturnError(issuer.UNAUTHORIZE)
	}

	assetN.Owner = newOwner
	assetJSON, err := json.Marshal(assetN)
	issuer.HandleError(err)
	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionNstdaStaff, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionNstdaStaff
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) GetAllNstdaStaff(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	var filterNstda = map[string]interface{}{}

	entityGetAll := entity.FilterGetAll{}
	interfaceNstda, err := issuer.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := interfaceNstda.(*entity.FilterGetAll)

	queryStringNstda, err := issuer.BuildQueryString(filterNstda)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryStringNstda)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, issuer.ReturnError(issuer.SKIPOVER)
	}

	arrNstda, err := core.FetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(arrNstda, func(i, j int) bool {
		return arrNstda[i].UpdatedAt.Before(arrNstda[j].UpdatedAt)
	})

	if len(arrNstda) == 0 {
		arrNstda = []*entity.TransectionReponse{}
	}

	return &entity.GetAllReponse{
		Data:  "All NstdaStaff",
		Obj:   arrNstda,
		Total: total,
	}, nil
}

func (s *SmartContract) FilterNstdaStaff(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionNstdaStaff, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assetNstda []*entity.TransectionNstdaStaff
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionNstdaStaff
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(queryResponse.Value, &m); err != nil {
			return nil, err
		}

		if val, ok := m[key]; ok && fmt.Sprintf("%v", val) == value {
			assetNstda = append(assetNstda, &asset)
		}
	}

	sort.Slice(assetNstda, func(i, j int) bool {
		return assetNstda[i].UpdatedAt.After(assetNstda[j].UpdatedAt)
	})

	return assetNstda, nil
}
