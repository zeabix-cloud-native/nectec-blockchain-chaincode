package entity

type TransectionFarmer struct {
	Id      string `json:"id"`
	CertId  string `json:"certId"`
	Owner   string `json:"owner"`
	OrgName string `json:"orgName"`
}

type Pagination struct {
	Skip  string `json:"skip"`
	Limit string `json:"limit"`
}
