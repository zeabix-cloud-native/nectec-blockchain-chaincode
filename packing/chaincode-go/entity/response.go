package entity

import "time"

type TransectionReponse struct {
	Id             string  `json:"id"`
	OrderID        string  `json:"orderId"`
	FarmerID       string  `json:"farmerId"`
	ForecastWeight float32 `json:"forecastWeight"`
	ActualWeight   float32 `json:"actualWeight"`
	// IsPackerSaved  bool      `json:"isPackerSaved"`
	SavedTime string `json:"savedTime"`
	// IsApproved     bool      `json:"isApproved"`
	ApprovedDate  string    `json:"approvedDate"`
	ApprovedType  string    `json:"approvedType"`
	FinalWeight   float32   `json:"finalWeight"`
	Remark        string    `json:"remark"`
	PackerId      string    `json:"packerId"`
	Gmp           string    `json:"gmp"`
	Gap           string    `json:"gap"`
	ProcessStatus int       `json:"processStatus"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}
