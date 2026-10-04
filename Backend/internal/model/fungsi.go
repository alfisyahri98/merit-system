package model

type Fungsi struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
}

func (Fungsi) TableName() string { return "fungsi" }
