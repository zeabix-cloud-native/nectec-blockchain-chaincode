package entity

type Transection struct {
	Id                  string `json:"id"`
	Prefix              string `json:"prefix"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	NationalID          string `json:"nationalID"`
	AddressRegistration string `json:"addressRegistration"`
	Address             string `json:"address"`
	VillageName         string `json:"villageName"`
	VillageNo           string `json:"villageNo"`
	Road                string `json:"road"`
	Alley               string `json:"alley"`
	Subdistrict         string `json:"subdistrict"`
	District            string `json:"district"`
	Province            string `json:"province"`
	ZipCode             string `json:"zipCode"`
	Phone               string `json:"phone"`
	MobilePhone         string `json:"mobilePhone"`
	Email               string `json:"email"`
	Owner               string `json:"owner"`
	OrgName             string `json:"orgName"`
}
