package entity

type TransectionReponse struct {
	Id     string `json:"id"`
	CertId string `json:"certId"`
	// Owner   string `json:"owner"`
	// OrgName string `json:"orgName"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}
