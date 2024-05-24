package core

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/regulator/chaincode-go/entity"
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

	getStringRegulator, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryRegulator, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringRegulator), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryRegulator.Close()

	var dataRegulator []*entity.TransectionReponse
	for queryRegulator.HasNext() {
		queryRes, err := queryRegulator.Next()
		if err != nil {
			return nil, err
		}

		var dataR entity.TransectionReponse
		err = json.Unmarshal(queryRes.Value, &dataR)
		if err != nil {
			return nil, err
		}

		dataRegulator = append(dataRegulator, &dataR)
	}

	return dataRegulator, nil
}
