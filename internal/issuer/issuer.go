package issuer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	UNAUTHORIZE string = "client is not authorized to delete this asset"
	TIMEFORMAT  string = "2006-01-02T15:04:05Z"
	SKIPOVER    string = "skip over total data"
)

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
