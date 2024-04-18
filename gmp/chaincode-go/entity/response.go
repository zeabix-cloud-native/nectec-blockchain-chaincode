package entity

import "time"

type TransectionReponse struct {
	Id                         string    `json:"id"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
	Address                    string    `json:"address"`
	UpdatedAt                  time.Time `json:"updatedAt"`
	CreatedAt                  time.Time `json:"createdAt"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}
