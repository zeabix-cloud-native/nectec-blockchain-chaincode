package farmer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/entity"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateAsset(
	ctx contractapi.TransactionContextInterface,
	args string,
) error {

	var input entity.TransectionFarmer

	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return fmt.Errorf("Unmarshal json string")
	}

	err := ctx.GetClientIdentity().AssertAttributeValue("farmer.creator", "true")
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

	formattedTime := time.Now().Format("2006-01-02T15:04:05Z")
	CreatedAt, _ := time.Parse("2006-01-02T15:04:05Z", formattedTime)

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

	var input entity.TransectionFarmer
	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return fmt.Errorf("Unmarshal json string")
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
		return fmt.Errorf("submitting client not authorized to update asset, does not own asset")
	}

	formattedTime := time.Now().Format("2006-01-02T15:04:05Z")
	UpdatedAt, _ := time.Parse("2006-01-02T15:04:05Z", formattedTime)

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
		return fmt.Errorf("submitting client not authorized to update asset, does not own asset")
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
		return fmt.Errorf("submitting client not authorized to update asset, does not own asset")
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
	log.Printf("Error creating farmer chaincode: %#c", asset)

	return &asset, nil
}

func (s *SmartContract) GetAllFarmer(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	orgName, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, err
	}

	var input entity.Pagination
	errInput := json.Unmarshal([]byte(args), &input)

	if errInput != nil {
		return nil, fmt.Errorf("Unmarshal json string")
	}

	limit := input.Limit
	skip := input.Skip

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	total := 0
	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		total++
	}

	selector := map[string]interface{}{
		"selector": map[string]interface{}{
			"orgName": orgName,
		},
		"limit": limit,
		"skip":  skip,
	}

	queryString, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryResults, _, err := ctx.GetStub().GetQueryResultWithPagination(string(queryString), int32(limit), "")
	if err != nil {
		return nil, err
	}
	defer queryResults.Close()

	var assets []*entity.TransectionReponse

	for queryResults.HasNext() {
		queryResponse, err := queryResults.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionReponse
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
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
		return "", fmt.Errorf("Failed to read clientID: %v", err)
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
		timestampStr := time.Unix(record.Timestamp.Seconds, int64(record.Timestamp.Nanos)).Format("2006-01-02T15:04:05Z")

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
