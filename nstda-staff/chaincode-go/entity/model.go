package entity

type TransectionNstdaStaff struct {
	Id      string `json:"id"`
	CertId  string `json:"certId"`
	Owner   string `json:"owner"`
	OrgName string `json:"orgName"`
}

type Pagination struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}
