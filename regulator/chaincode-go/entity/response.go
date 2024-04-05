package entity

type TransectionReponse struct {
	Id     string `json:"id"`
	CertId string `json:"certId"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}
