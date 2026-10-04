package model

type Pangkat struct {
	ID       int    `json:"id"`
	Jenis    string `json:"jenis"`
	Kode     string `json:"kode"`
	Nama     string `json:"nama"`
	Kelompok string `json:"kelompok"`
	Urutan   int    `json:"urutan"`
}

func (Pangkat) TableName() string { return "pangkat" }
