package entity

import "time"

type TransectionReponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetAllReponse struct {
	Data  string                `json:"data"`
	Obj   []*TransectionReponse `json:"obj"`
	Total int                   `json:"total"`
}

type TransactionHistory struct {
	TxId      string                `json:"tx_id"`
	IsDelete  bool                  `json:"isDelete"`
	Value     []*TransectionReponse `json:"value"`
	Timestamp string                `json:"timestamp"`
}
