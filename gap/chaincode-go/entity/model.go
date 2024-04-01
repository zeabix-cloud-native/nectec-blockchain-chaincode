package entity

type TransectionGAP struct {
	Id          string `json:"id"`
	CertID      string `json:"certId"`
	AreaCode    string `json:"areaCode"`
	AreaSize    string `json:"areaSize"`
	AreaStatus  string `json:"areaStatus"`
	OldAreaCode string `json:"oldAreaCode"`
	IssueDate   string `json:"issueDate"`
	ExpireDate  string `json:"expireDate"`
	District    string `json:"district"`
	Province    string `json:"province"`
	UpdatedDate string `json:"updatedDate"`
	Source      string `json:"source"`
	Owner       string `json:"owner"`
	OrgName     string `json:"orgName"`
}
