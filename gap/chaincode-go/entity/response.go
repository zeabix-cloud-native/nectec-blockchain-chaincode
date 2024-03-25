package entity

type TransectionReponse struct {
	Prefix      string `json:"prefix"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	NationalID  string `json:"nationalID"`
	Phone       string `json:"phone"`
	MobilePhone string `json:"mobilePhone"`
	Owner       string `json:"owner"`
	OrgName     string `json:"orgName"`
}
