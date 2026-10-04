package model

type Nivelering struct {
	ID     int    `json:"id"`
	Kode   string `json:"kode"`
	Urutan int    `json:"urutan"`
}

func (Nivelering) TableName() string { return "nivelering" }
