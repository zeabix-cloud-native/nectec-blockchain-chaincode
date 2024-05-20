package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gap/chaincode-go/entity"
)

func UnmarshalGap(args string) (entity.TransectionGAP, error) {
	var input entity.TransectionGAP
	err := json.Unmarshal([]byte(args), &input)
	if err != nil {
		return entity.TransectionGAP{}, fmt.Errorf("unmarshal json string: %v", err)
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
func SetFilter(input entity.FilterGetAll) map[string]interface{} {
	var filter = map[string]interface{}{}
	if input.CertID != nil {
		filter["certId"] = *input.CertID
	}
	if input.AreaCode != nil {
		filter["areaCode"] = *input.AreaCode
	}
	if input.Province != nil {
		filter["province"] = *input.Province
	}
	if input.District != nil {
		filter["district"] = *input.District
	}
	if input.AreaRaiFrom != nil && input.AreaRaiTo != nil {
		filter["areaRai"] = map[string]interface{}{
			"$gte": *input.AreaRaiFrom,
			"$lte": *input.AreaRaiTo,
		}
	}
	if input.IssueDate != nil {
		filter["issueDate"] = *input.IssueDate
	}
	if input.ExpireDate != nil {
		filter["expireDate"] = *input.ExpireDate
	}

	if input.AvailableGap != nil {
		filter["farmerId"] = ""
	}

	return filter
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

func FetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input entity.FilterGetAll, filter map[string]interface{}) ([]*entity.TransectionReponse, error) {

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

func GetTimeNow() time.Time {
	formattedTime := time.Now().Format(entity.TimeFormat)
	CreatedAt, _ := time.Parse(entity.TimeFormat, formattedTime)
	return CreatedAt
}
