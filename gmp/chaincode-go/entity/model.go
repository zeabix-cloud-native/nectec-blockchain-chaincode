package entity

type TransectionGMP struct {
	Id                         string `json:"id"`
	PackingHouseRegisterNumber string `json:"packingHouseRegisterNumber"`
	Address                    string `json:"address"`
	Owner                      string `json:"owner"`
	OrgName                    string `json:"orgName"`
}

type Pagination struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}
