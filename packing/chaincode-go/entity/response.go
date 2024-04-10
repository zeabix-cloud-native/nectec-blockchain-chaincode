package entity

type TransectionReponse struct {
	Id             string `json:"id"`
	OrderID        string `json:"orderId"`
	FarmerID       string `json:"farmerId"`
	ForecastWeight string `json:"forecastWeight"`
	ActualWeight   string `json:"actualWeight"`
	IsPackerSaved  bool   `json:"isPackerSaved"`
	SavedTime      string `json:"savedTime"`
	IsApproved     bool   `json:"isApproved"`
	ApprovedDate   string `json:"approvedDate"`
	ApprovedType   string `json:"approvedType"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}