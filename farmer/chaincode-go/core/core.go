package core

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/entity"
)

func UnmarshalFarmer(args string) (entity.TransectionFarmer, error) {
	var input entity.TransectionFarmer
	err := json.Unmarshal([]byte(args), &input)
	if err != nil {
		return entity.TransectionFarmer{}, fmt.Errorf("unmarshal json string: %v", err)
	}
	return input, nil
}

func UnmarshalGetAll(args string) (entity.FilterGetAll, error) {
	var input entity.FilterGetAll
	err := json.Unmarshal([]byte(args), &input)
	if err != nil {
		return entity.FilterGetAll{}, fmt.Errorf("unmarshal json string: %v", err)
	}
	return input, nil
}

func BuildQueryString(filter map[string]interface{}) (string, error) {
	selector := map[string]interface{}{
		"selector": filter,
	}
	queryString, err := json.Marshal(selector)
	if err != nil {
		return "", err
	}
	return string(queryString), nil
}

func CountTotalResults(ctx contractapi.TransactionContextInterface, queryString string) (int, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return 0, err
	}
	defer resultsIterator.Close()

	total := 0
	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return 0, err
		}
		total++
	}
	return total, nil
}

func FetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input entity.FilterGetAll) ([]*entity.TransectionReponse, error) {
	var filter = map[string]interface{}{}

	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	queryString, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryResults, _, err := ctx.GetStub().GetQueryResultWithPagination(string(queryString), int32(input.Limit), "")
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

	return assets, nil
}
