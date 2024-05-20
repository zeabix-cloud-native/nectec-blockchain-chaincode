package entity

import "time"

const (
	UNAUTHORIZE  string = "client is not authorized to delete this asset"
	UNAUTHORIZE1 string = "client is not authorized to delete this asset"
)

type TransectionFarmer struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type FilterGetAll struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}
