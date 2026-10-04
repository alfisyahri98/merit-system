package model

import (
	"time"

	"gorm.io/gorm"
)

// MasterPendidikan: daftar jenis pendidikan/pelatihan.
// Jenis: DIKPOL, DIKUM, DIKBANG, PELATIHAN.
type MasterPendidikan struct {
	ID     int    `json:"id"`
	Jenis  string `json:"jenis"`
	Kode   string `json:"kode"`
	Nama   string `json:"nama"`
	Urutan *int   `json:"urutan"`
}

func (MasterPendidikan) TableName() string { return "master_pendidikan" }

// Kualifikasi: pendidikan / pelatihan / sertifikasi yang pernah diikuti personel.
type Kualifikasi struct {
	ID                 int64          `json:"id"`
	PersonelID         int64          `json:"personel_id"`
	MasterPendidikanID int            `json:"master_pendidikan_id"`
	Institusi          *string        `json:"institusi"`
	Jurusan            *string        `json:"jurusan"`
	TahunLulus         int            `json:"tahun_lulus"`
	NomorDokumen       *string        `json:"nomor_dokumen"`
	Keterangan         *string        `json:"keterangan"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-"`

	Pendidikan *MasterPendidikan `gorm:"foreignKey:MasterPendidikanID" json:"pendidikan,omitempty"`
}

func (Kualifikasi) TableName() string { return "kualifikasi" }
