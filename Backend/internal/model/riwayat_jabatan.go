package model

import (
	"time"

	"gorm.io/gorm"
)

type RiwayatJabatan struct {
	ID            int64          `json:"id"`
	PersonelID    int64          `json:"personel_id"`
	NamaJabatan   string         `json:"nama_jabatan"`
	SatkerID      int            `json:"satker_id"`
	FungsiID      *int           `json:"fungsi_id"`
	NiveleringID  *int           `json:"nivelering_id"`
	TmtMulai      Date           `json:"tmt_mulai"`
	TmtSelesai    *Date          `json:"tmt_selesai"`
	StatusJabatan string         `json:"status_jabatan"`
	NomorSkep     *string        `json:"nomor_skep"`
	Keterangan    *string        `json:"keterangan"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-"`

	Satker     *Satker     `gorm:"foreignKey:SatkerID" json:"satker,omitempty"`
	Fungsi     *Fungsi     `gorm:"foreignKey:FungsiID" json:"fungsi,omitempty"`
	Nivelering *Nivelering `gorm:"foreignKey:NiveleringID" json:"nivelering,omitempty"`
}

func (RiwayatJabatan) TableName() string {
	return "riwayat_jabatan"
}
