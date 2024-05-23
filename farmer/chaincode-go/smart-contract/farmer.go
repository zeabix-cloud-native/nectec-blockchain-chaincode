package farmer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/entity"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateFarmer(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	entityFarmer := entity.TransectionFarmer{}
	inputInterface, err := issuer.Unmarshal(args, entityFarmer)
	if err != nil {
		return err
	}
	input := inputInterface.(*entity.TransectionFarmer)

	// err := ctx.GetClientIdentity().AssertAttributeValue("farmer.creator", "true")
	// if err != nil {
	// 	return fmt.Errorf("submitting client not authorized to create asset, does not have abac.creator role")
	// }

	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("submitting client not authorized to create asset, does not have farmer.creator role")
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

	CreatedAt := issuer.GetTimeNow()

	asset := entity.TransectionFarmer{
		Id:        input.Id,
		CertId:    input.CertId,
		Owner:     clientID,
		OrgName:   orgName,
		UpdatedAt: CreatedAt,
		CreatedAt: CreatedAt,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(input.Id, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface,
	args string) error {
	entityType := entity.TransectionFarmer{}
	inputInterface, err := issuer.Unmarshal(args, entityType)
	if err != nil {
		return err
	}
	input := inputInterface.(*entity.TransectionFarmer)

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

	UpdatedAt := issuer.GetTimeNow()

	asset.Id = input.Id
	asset.CertId = input.CertId
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

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*entity.TransectionFarmer, error) {

	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset entity.TransectionFarmer
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) GetAllFarmer(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

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
		Data:  "All Farmer",
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

func (s *SmartContract) FilterFarmer(ctx contractapi.TransactionContextInterface, key, value string) ([]*entity.TransectionFarmer, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*entity.TransectionFarmer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionFarmer
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

func (s *SmartContract) GetHistoryForKey(ctx contractapi.TransactionContextInterface, key string) ([]*entity.TransactionHistory, error) {
	// Get the history for the specified key
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for key %s: %v", key, err)
	}
	defer resultsIterator.Close()

	var history []*entity.TransactionHistory
	var assetsValue []*entity.TransectionReponse
	for resultsIterator.HasNext() {
		// Get the next history record
		record, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next history record for key %s: %v", key, err)
		}

		var asset entity.TransectionReponse
		if !record.IsDelete {
			err = json.Unmarshal(record.Value, &asset)
			if err != nil {
				return nil, err
			}
			assetsValue = append(assetsValue, &asset)

		} else {
			assetsValue = []*entity.TransectionReponse{}
		}
		// Convert the timestamp to string in the desired format
		timestampStr := time.Unix(record.Timestamp.Seconds, int64(record.Timestamp.Nanos)).Format(issuer.TIMEFORMAT)

		historyRecord := &entity.TransactionHistory{
			TxId:      record.TxId,
			Value:     assetsValue,
			Timestamp: timestampStr,
			IsDelete:  record.IsDelete,
		}

		history = append(history, historyRecord)
	}

	return history, nil
}

func (s *SmartContract) GetLastIdFarmer(ctx contractapi.TransactionContextInterface) string {
	// Query to get all records sorted by ID in descending order
	query := `{
			"selector": {},
			"sort": [{"_id": "desc"}],
			"limit": 1,
			"use_index": "index-id"
	}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return "error querying CouchDB"
	}
	defer resultsIterator.Close()

	// Check if there is a result
	if !resultsIterator.HasNext() {
		return ""
	}

	// Get the first (and only) result
	queryResponse, err := resultsIterator.Next()
	if err != nil {
		return "error iterating query results"
	}

	var result struct {
		Id string `json:"id"`
	}

	// Unmarshal the result into the result struct
	if err := json.Unmarshal(queryResponse.Value, &result); err != nil {
		return "error unmarshalling document"
	}

	return result.Id
}

func (s *SmartContract) SaveUserEvent(ctx contractapi.TransactionContextInterface, args string) {
	assetJSON, _ := json.Marshal(args)
	ctx.GetStub().SetEvent("SaveUserEvent", assetJSON)
}

func (s *SmartContract) CreateFarmerCsv(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {
	var inputs []entity.TransectionFarmer
	var eventPayloads []entity.TransectionFarmer

	errInput := json.Unmarshal([]byte(args), &inputs)
	if errInput != nil {
		return fmt.Errorf("failed to unmarshal JSON array: %v", errInput)
	}

	for _, input := range inputs {
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

		asset := entity.TransectionFarmer{
			Id:        input.Id,
			CertId:    input.CertId,
			Owner:     clientID,
			OrgName:   orgName,
			UpdatedAt: input.CreatedAt,
			CreatedAt: input.UpdatedAt,
		}

		assetJSON, err := json.Marshal(asset)
		eventPayloads = append(eventPayloads, asset)
		if err != nil {
			return fmt.Errorf("failed to marshal asset JSON: %v", err)
		}

		err = ctx.GetStub().PutState(input.Id, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put state for asset %s: %v", input.Id, err)
		}

		fmt.Printf("Asset %s created successfully\n", input.Id)

	}

	eventPayloadJSON, err := json.Marshal(eventPayloads)
	if err != nil {
		return fmt.Errorf("failed to marshal asset JSON: %v", err)
	}
	ctx.GetStub().SetEvent("batchCreatedUserEvent", eventPayloadJSON)

	return nil
}
