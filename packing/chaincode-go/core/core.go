package core

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/packing/chaincode-go/entity"
)

func SetFilter(input *entity.FilterGetAll) map[string]interface{} {
	var filter = map[string]interface{}{}

	if input.Gap != nil {
		filter["gap"] = *input.Gap
	}

	if input.CertID != nil {
		filter["certId"] = *input.CertID
	}

	if input.Search != nil {
		filter["search"] = *input.Search
	}

	if (input.FarmerID != nil) {
		filter["farmerId"] = *input.FarmerID
	}

	if input.StartDate != nil && input.EndDate != nil {
		filter["createdAt"] = map[string]interface{}{
			"$gte": *input.StartDate,
			"$lte": *input.EndDate,
		}
	}

	if input.ForecastWeightFrom != nil && input.ForecastWeightTo != nil {
		filter["forecastWeight"] = map[string]interface{}{
			"$gte": *input.ForecastWeightFrom,
			"$lte": *input.ForecastWeightTo,
		}
	}

	if input.ProcessStatus != nil {
		filter["processStatus"] = *input.ProcessStatus
	}

	return filter
}

func FetchResultsWithPagination(ctx contractapi.TransactionContextInterface, input *entity.FilterGetAll, filter map[string]interface{}) ([]*entity.TransectionReponse, error) {
	search, searchExists := filter["search"]

	if searchExists {
		delete(filter, "search")
	}

	// Initialize the base selector
	selector := map[string]interface{}{
		"selector": filter,
	}

	if searchExists && search != "" {
		selector["selector"] = map[string]interface{}{
			"$and": []map[string]interface{}{
				filter,
				{
					"$or": []map[string]interface{}{
						{"gmp": map[string]interface{}{"$regex": search}},
						{"packingHouseName": map[string]interface{}{"$regex": search}},
					},
				},
			},
		}
	}

	if input.Skip != 0 || input.Limit != 0 {
		selector["skip"] = input.Skip
		selector["limit"] = input.Limit
	}

	getStringPacking, err := json.Marshal(selector)
	if err != nil {
		return nil, err
	}

	queryPacking, _, err := ctx.GetStub().GetQueryResultWithPagination(string(getStringPacking), int32(input.Limit), "")
	if err != nil {
		return nil, err
	}
	defer queryPacking.Close()

	var dataPacking []*entity.TransectionReponse
	for queryPacking.HasNext() {
		queryResponse, err := queryPacking.Next()
		if err != nil {
			return nil, err
		}

		var asset entity.TransectionReponse
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		dataPacking = append(dataPacking, &asset)
	}

	return dataPacking, nil
}
