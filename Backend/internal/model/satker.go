package model

type Satker struct {
	ID       int    `json:"id"`
	Kode     string `json:"kode"`
	Nama     string `json:"nama"`
	ParentID *int   `json:"parent_id"`
	Tingkat  string `json:"tingkat"`
}

func (Satker) TableName() string {
	return "satker"
}
