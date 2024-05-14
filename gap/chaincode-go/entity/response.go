package entity

import "time"

type TransectionReponse struct {
	Id          string    `json:"id"`
	CertID      string    `json:"certId"`
	AreaCode    string    `json:"areaCode"`
	AreaRai     float32   `json:"areaRai"`
	AreaStatus  string    `json:"areaStatus"`
	OldAreaCode string    `json:"oldAreaCode"`
	IssueDate   string    `json:"issueDate"`
	ExpireDate  string    `json:"expireDate"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	UpdatedDate string    `json:"updatedDate"`
	Source      string    `json:"source"`
	FarmerID    string    `json:"farmerId"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}

type GetByCertIDReponse struct {
	Data string              `json:"data"`
	Obj  *TransectionReponse `json:"obj"`
}
