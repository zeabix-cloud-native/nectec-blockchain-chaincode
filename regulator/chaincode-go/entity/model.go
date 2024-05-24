package entity

import "time"

type TransectionRegulator struct {
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
