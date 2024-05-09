package entity

import "time"

type TransectionPacking struct {
	Id             string    `json:"id"`
	OrderID        string    `json:"orderId"`
	FarmerID       string    `json:"farmerId"`
	ForecastWeight float32   `json:"forecastWeight"`
	ActualWeight   float32   `json:"actualWeight"`
	SavedTime      string    `json:"savedTime"`
	ApprovedDate   string    `json:"approvedDate"`
	ApprovedType   string    `json:"approvedType"`
	FinalWeight    float32   `json:"finalWeight"`
	Remark         string    `json:"remark"`
	PackerId       string    `json:"packerId"`
	Gmp            string    `json:"gmp"`
	Gap            string    `json:"gap"` // รหัสซื้อขาย
	ProcessStatus  int       `json:"processStatus"`
	Owner          string    `json:"owner"`
	OrgName        string    `json:"orgName"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

type FilterGetAll struct {
	Skip               int      `json:"skip"`
	Limit              int      `json:"limit"`
	PackerId           *string  `json:"packerId"`
	Gap                *string  `json:"gap"`
	StartDate          *string  `json:"startDate"`
	EndDate            *string  `json:"endDate"`
	ForecastWeightFrom *float32 `json:"forecastWeightFrom"`
	ForecastWeightTo   *float32 `json:"forecastWeightTo"`
	ProcessStatus      *int     `json:"processStatus"`
}
