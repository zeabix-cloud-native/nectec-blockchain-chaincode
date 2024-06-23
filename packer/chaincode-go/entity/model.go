package entity

import "time"

type TransectionPacker struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UserId    string    `json:"userId"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	PackerGmp PackerGmp `json:"packerGmp"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type FilterGetAll struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	PackerGmp string `json:"packerGmp"`
}

type PackerGmp struct {
	Id                         string    `json:"id"`
	PackerId 				   string    `json:"packerId"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
	Address                    string    `json:"address"`
	PackingHouseName           string    `json:"packingHouseName"`
	UpdatedDate                string    `json:"updatedDate"`
	Source                     string    `json:"source"`
	Owner                      string    `json:"owner"`
	OrgName                    string    `json:"orgName"`
	UpdatedAt                  time.Time `json:"updatedAt"`
	CreatedAt                  time.Time `json:"createdAt"`
}