package core

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packer/chaincode-go/entity"
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

	getStringPacker, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryPacker, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringPacker), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryPacker.Close()

	var dataPacker []*entity.TransectionReponse
	for queryPacker.HasNext() {
		queryRes, err := queryPacker.Next()
		if err != nil {
			return nil, err
		}

		var dataP entity.TransectionReponse
		err = json.Unmarshal(queryRes.Value, &dataP)
		if err != nil {
			return nil, err
		}

		dataPacker = append(dataPacker, &dataP)
	}

	return dataPacker, nil
}
