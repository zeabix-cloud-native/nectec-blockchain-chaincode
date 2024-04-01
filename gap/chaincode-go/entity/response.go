package entity

import "time"

type TransectionReponse struct {
	CertID      string    `json:"certId"`
	AreaCode    string    `json:"areaCode"`
	AreaSize    string    `json:"areaSize"`
	AreaStatus  string    `json:"areaStatus"`
	OldAreaCode string    `json:"oldAreaCode"`
	IssueDate   time.Time `json:"issueDate"`
	ExpireDate  time.Time `json:"expireDate"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	UpdatedDate time.Time `json:"updatedDate"`
	Source      string    `json:"source"`
	Owner       string    `json:"owner"`
	OrgName     string    `json:"orgName"`
}
