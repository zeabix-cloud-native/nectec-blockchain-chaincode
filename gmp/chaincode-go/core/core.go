package core

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/gmp/chaincode-go/entity"
)

func SetFilter(input *entity.FilterGetAll) map[string]interface{} {
	var filter = map[string]interface{}{}

	if input.PackingHouseRegisterNumber != nil {
		filter["packingHouseRegisterNumber"] = input.PackingHouseRegisterNumber
	}

	if input.Address != nil {
		filter["address"] = input.Address
	}

	return filter
}

func FetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *entity.FilterGetAll, filter map[string]interface{}) ([]*entity.TransectionReponse, error) {
	selector := map[string]interface{}{
		"selector": filter,
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getStringGmp, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryGmp, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringGmp), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryGmp.Close()

	var dataGmp []*entity.TransectionReponse
	for queryGmp.HasNext() {
		queryRes, err := queryGmp.Next()
		if err != nil {
			return nil, err
		}

		var dataG entity.TransectionReponse
		err = json.Unmarshal(queryRes.Value, &dataG)
		if err != nil {
			return nil, err
		}

		dataGmp = append(dataGmp, &dataG)
	}

	return dataGmp, nil
}
