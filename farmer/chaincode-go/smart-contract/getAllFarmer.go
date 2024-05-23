package farmer

import (
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/core"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/farmer/chaincode-go/entity"
	"github.com/zeabix-cloud-native/nstda-blockchain-chaincode/internal/issuer"
)

func (s *SmartContract) GetAllFarmer(ctx contractapi.TransactionContextInterface, args string) (*entity.GetAllReponse, error) {

	var filter = map[string]interface{}{}

	entityGetAll := entity.FilterGetAll{}
	inputInterface, err := issuer.Unmarshal(args, entityGetAll)
	if err != nil {
		return nil, err
	}
	input := inputInterface.(*entity.FilterGetAll)

	queryString, err := issuer.BuildQueryString(filter)
	if err != nil {
		return nil, err
	}

	total, err := issuer.CountTotalResults(ctx, queryString)
	if err != nil {
		return nil, err
	}

	if input.Skip > total {
		return nil, fmt.Errorf(issuer.SKIPOVER)
	}

	assets, err := core.FetchResultsWithPagination(ctx, input)
	if err != nil {
		return nil, err
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].UpdatedAt.Before(assets[j].UpdatedAt)
	})

	if len(assets) == 0 {
		assets = []*entity.TransectionReponse{}
	}

	return &entity.GetAllReponse{
		Data:  "All Farmer",
		Obj:   assets,
		Total: total,
	}, nil
}
