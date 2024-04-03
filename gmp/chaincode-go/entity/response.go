package entity

type TransectionReponse struct {
	Id                         string `json:"id"`
	PackingHouseRegisterNumber string `json:"packingHouseRegisterNumber"`
	Address                    string `json:"address"`
	// Owner                      string `json:"owner"`
	// OrgName                    string `json:"orgName"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}
