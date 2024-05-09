package entity

import "time"

type TransectionGMP struct {
	Id                         string    `json:"id"`
	Name                       string    `json:"name"`
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

type Pagination struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type FilterGetAll struct {
	Skip                       int     `json:"skip"`
	Limit                      int     `json:"limit"`
	Name                       *string `json:"name"`
	PackingHouseRegisterNumber *string `json:"packingHouseRegisterNumber"`
	Address                    *string `json:"address"`
}
