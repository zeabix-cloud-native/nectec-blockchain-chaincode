package core

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/exporter/chaincode-go/entity"
)

func FetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *entity.FilterGetAll) ([]*entity.TransectionReponse, error) {
	var filter = map[string]interface{}{}

	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getString, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryExporter, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getString), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryExporter.Close()

	var dataExporters []*entity.TransectionReponse
	for queryExporter.HasNext() {
		queryRes, err := queryExporter.Next()
		if err != nil {
			return nil, err
		}

		var dataF entity.TransectionReponse
		err = json.Unmarshal(queryRes.Value, &dataF)
		if err != nil {
			return nil, err
		}

		dataExporters = append(dataExporters, &dataF)
	}

	return dataExporters, nil
}
