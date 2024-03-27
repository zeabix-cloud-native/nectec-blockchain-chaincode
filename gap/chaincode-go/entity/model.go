package entity

type TransectionGAP struct {
	Id                string         `json:"id"`
	AgriStandard      string         `json:"agriStandard"`
	FarmOwner         Owner          `json:"farmOwner"`
	FarmOwnerJuristic JuristicPerson `json:"farmOwnerJuristicPerson"`
	FarmLocation      Location       `json:"farmLocation"`
	RegisterPlants    []Plant        `json:"registerPlants"`
	Owner             string         `json:"owner"`
	OrgName           string         `json:"orgName"`
}

type Owner struct {
	Prefix                string `json:"prefix"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	IdCard                string `json:"idCard"`
	HouseRegistrationCode string `json:"houseRegistrationCode"`
	HouseNo               string `json:"houseNo"`
	Road                  string `json:"road"`
	SubDistrict           string `json:"subDistrict"`
	District              string `json:"district"`
	Province              string `json:"province"`
	PostalCode            string `json:"postalCode"`
	Phone                 string `json:"phone"`
	MobilePhone           string `json:"mobilePhone"`
	Email                 string `json:"email"`
}

type JuristicPerson struct {
	JuristicId             string `json:"juristicId"`
	LegalEntityRegisNumber string `json:"legalEntityRegisNumber"`
	Prefix                 string `json:"prefix"`
	FirstName              string `json:"firstName"`
	LastName               string `json:"lastName"`
	IdCard                 string `json:"idCard"`
	HouseRegistrationCode  string `json:"houseRegistrationCode"`
	HouseNo                string `json:"houseNo"`
	Road                   string `json:"road"`
	SubDistrict            string `json:"subDistrict"`
	District               string `json:"district"`
	Province               string `json:"province"`
	PostalCode             string `json:"postalCode"`
	Phone                  string `json:"phone"`
	MobilePhone            string `json:"mobilePhone"`
	Email                  string `json:"email"`
}

type Location struct {
	VillageName   string `json:"villageName"`
	Road          string `json:"road"`
	SubDistrict   string `json:"subDistrict"`
	District      string `json:"district"`
	Province      string `json:"province"`
	CertifiedArea int    `json:"certifiedArea"`
}

type Plant struct {
	PlantType            string `json:"plantType"`
	Area                 int    `json:"area"`
	PlantAge             int    `json:"plantAge"`
	ProductionPeriod     int    `json:"productionPeriod"`
	ExpectedHarvestTime  int    `json:"expectedHarvestTime"`
	TotalProductionYear  int    `json:"totalProductionPerYear"`
	IdentificationNumber string `json:"identificationNumber"`
}
