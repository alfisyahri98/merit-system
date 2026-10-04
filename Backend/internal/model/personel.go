package model

import (
	"time"

	"gorm.io/gorm"
)

type Personel struct {
	ID           int64          `json:"id"`
	Jenis        string         `json:"jenis"`
	NrpNip       string         `json:"nrp_nip"`
	Nama         string         `json:"nama"`
	PangkatID    int            `json:"pangkat_id"`
	SatkerID     int            `json:"satker_id"`
	TempatLahir  *string        `json:"tempat_lahir"`
	TanggalLahir Date           `json:"tanggal_lahir"`
	Status       string         `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-"`

	Pangkat        *Pangkat         `gorm:"foreignKey:PangkatID" json:"pangkat,omitempty"`
	Satker         *Satker          `gorm:"foreignKey:SatkerID" json:"satker,omitempty"`
	RiwayatJabatan []RiwayatJabatan `gorm:"foreignKey:PersonelID" json:"riwayat_jabatan,omitempty"`
	JabatanSaatIni *RiwayatJabatan  `gorm:"-" json:"jabatan_saat_ini,omitempty"`
	Sertifikasi    []Kualifikasi    `gorm:"foreignKey:PersonelID" json:"sertifikasi,omitempty"` // kualifikasi jenis PELATIHAN
}

func (Personel) TableName() string { return "personel" }
