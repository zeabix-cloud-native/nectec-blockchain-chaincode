package issuer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	UNAUTHORIZE   string = "client is not authorized this asset"
	TIMEFORMAT    string = "2006-01-02T15:04:05Z"
	SKIPOVER      string = "skip over total data"
	DATAUNMARSHAL string = "unmarshal json string"
)

type GetAllType struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

func Unmarshal(args string, entityType interface{}) (interface{}, error) {
	entityValue := reflect.New(reflect.TypeOf(entityType)).Interface()
	err := json.Unmarshal([]byte(args), entityValue)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json string: %v", err)
	}
	return entityValue, nil
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

func GetTimeNow() time.Time {
	formattedTime := time.Now().Format(TIMEFORMAT)
	CreatedAt, _ := time.Parse(TIMEFORMAT, formattedTime)
	return CreatedAt
}

func GetAllNotFilter(ctx contractapi.TransactionContextInterface, input GetAllType, resultType interface{}) ([]interface{}, error) {
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

	var results []interface{}
	resultTypeValue := reflect.TypeOf(resultType).Elem()

	for queryResults.HasNext() {
		queryResponse, err := queryResults.Next()
		if err != nil {
			return nil, err
		}
		resultInstance := reflect.New(resultTypeValue).Interface()

		err = json.Unmarshal(queryResponse.Value, &resultInstance)
		if err != nil {
			return nil, err
		}

		results = append(results, resultInstance)
	}

	return results, nil
}

func ReturnError(data string) error {
	return fmt.Errorf(data)
}
