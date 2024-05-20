package entity

import "time"

const (
	UNAUTHORIZE string = "client is not authorized to delete this asset"
	TimeFormat  string = "2006-01-02T15:04:05Z"
	SkipOver    string = "skip over total data"
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
