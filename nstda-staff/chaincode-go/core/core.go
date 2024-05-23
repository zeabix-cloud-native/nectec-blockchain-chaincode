package core

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/nstda-staff/chaincode-go/entity"
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

	getStringNstda, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryNstda, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringNstda), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryNstda.Close()

	var dataNstda []*entity.TransectionReponse
	for queryNstda.HasNext() {
		queryRes, err := queryNstda.Next()
		if err != nil {
			return nil, err
		}

		var dataP entity.TransectionReponse
		err = json.Unmarshal(queryRes.Value, &dataP)
		if err != nil {
			return nil, err
		}

		dataNstda = append(dataNstda, &dataP)
	}

	return dataNstda, nil
}
