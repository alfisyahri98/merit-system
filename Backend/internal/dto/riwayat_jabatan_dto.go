package dto

import (
	"Backend/internal/model"
	"strings"
)

// RiwayatJabatanRequest: body untuk POST & PUT riwayat jabatan.
type RiwayatJabatanRequest struct {
	NamaJabatan   string      `json:"nama_jabatan" binding:"required,max=300"`
	SatkerID      int         `json:"satker_id" binding:"required,gt=0"`
	FungsiID      *int        `json:"fungsi_id" binding:"omitempty,gt=0"`
	NiveleringID  *int        `json:"nivelering_id" binding:"omitempty,gt=0"`
	TmtMulai      model.Date  `json:"tmt_mulai"`
	TmtSelesai    *model.Date `json:"tmt_selesai"` // null = jabatan saat ini
	StatusJabatan string      `json:"status_jabatan" binding:"omitempty,oneof=DEFINITIF PLT PLH"`
	NomorSkep     *string     `json:"nomor_skep" binding:"omitempty,max=100"`
	Keterangan    *string     `json:"keterangan" binding:"omitempty,max=2000"`
}

// Validate: aturan bisnis. tanggalLahir dipakai untuk cek TMT masuk akal.
func (r *RiwayatJabatanRequest) Validate(tanggalLahir model.Date) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(r.NamaJabatan) == "" {
		errs["nama_jabatan"] = "nama jabatan wajib diisi"
	}

	if r.TmtMulai.IsZero() {
		errs["tmt_mulai"] = "tanggal mulai menjabat wajib diisi (YYYY-MM-DD)"
	} else if !tanggalLahir.IsZero() && !r.TmtMulai.After(tanggalLahir.Time) {
		errs["tmt_mulai"] = "tanggal mulai menjabat harus setelah tanggal lahir personel"
	}

	if r.TmtSelesai != nil && !r.TmtMulai.IsZero() && r.TmtSelesai.Before(r.TmtMulai.Time) {
		errs["tmt_selesai"] = "tanggal berakhir tidak boleh sebelum tanggal mulai"
	}
	return errs
}

// ApplyTo: salin isi request ke model.
func (r *RiwayatJabatanRequest) ApplyTo(rj *model.RiwayatJabatan) {
	rj.NamaJabatan = strings.TrimSpace(r.NamaJabatan)
	rj.SatkerID = r.SatkerID
	rj.FungsiID = r.FungsiID
	rj.NiveleringID = r.NiveleringID
	rj.TmtMulai = r.TmtMulai
	rj.TmtSelesai = r.TmtSelesai
	rj.StatusJabatan = r.StatusJabatan
	if rj.StatusJabatan == "" {
		rj.StatusJabatan = "DEFINITIF"
	}
	rj.NomorSkep = r.NomorSkep
	rj.Keterangan = r.Keterangan
}
