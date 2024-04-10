package entity

type TransectionReponse struct {
	Id          string `json:"id"`
	CertID      string `json:"certId"`
	AreaCode    string `json:"areaCode"`
	AreaRai     string `json:"areaRai"`
	AreaStatus  string `json:"areaStatus"`
	OldAreaCode string `json:"oldAreaCode"`
	IssueDate   string `json:"issueDate"`
	ExpireDate  string `json:"expireDate"`
	District    string `json:"district"`
	Province    string `json:"province"`
	UpdatedDate string `json:"updatedDate"`
	Source      string `json:"source"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}
