package entity

type TransectionPacking struct {
	Id             string `json:"id"`
	OrderID        string `json:"orderId"`
	FarmerID       string `json:"farmerId"`
	ForecastWeight string `json:"forecastWeight"`
	ActualWeight   string `json:"actualWeight"`
	IsPackerSaved  bool   `json:"isPackerSaved"`
	SavedTime      string `json:"savedTime"`
	IsApproved     bool   `json:"isApproved"` // update status
	ApprovedDate   string `json:"approvedDate"`
	ApprovedType   string `json:"approvedType"`
	Owner          string `json:"owner"`
	OrgName        string `json:"orgName"`
}

type Pagination struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}
