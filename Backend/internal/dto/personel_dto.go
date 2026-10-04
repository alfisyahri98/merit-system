package dto

import (
	"Backend/internal/model"
	"strings"
	"time"
)

// PersonelRequest: body untuk POST & PUT /personel.
// Lapis validasi 1 (tag binding): tipe, wajib isi, nilai yang diizinkan.
type PersonelRequest struct {
	Jenis        string     `json:"jenis" binding:"required,oneof=POLRI PNS"`
	NrpNip       string     `json:"nrp_nip" binding:"required,numeric"`
	Nama         string     `json:"nama" binding:"required,max=150"`
	PangkatID    int        `json:"pangkat_id" binding:"required,gt=0"`
	SatkerID     int        `json:"satker_id" binding:"required,gt=0"`
	TempatLahir  *string    `json:"tempat_lahir" binding:"omitempty,max=100"`
	TanggalLahir model.Date `json:"tanggal_lahir"`
	Status       string     `json:"status" binding:"omitempty,oneof=AKTIF PENSIUN MENINGGAL DIBERHENTIKAN MUTASI_KELUAR"`
}

// Validate: lapis validasi 2 — aturan bisnis yang tidak bisa ditulis di tag.
// Return map field → pesan error (kosong = valid).
func (r *PersonelRequest) Validate() map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(r.Nama) == "" {
		errs["nama"] = "nama wajib diisi"
	}
	if r.Jenis == "POLRI" && len(r.NrpNip) != 8 {
		errs["nrp_nip"] = "NRP POLRI harus 8 digit angka"
	}
	if r.Jenis == "PNS" && len(r.NrpNip) != 18 {
		errs["nrp_nip"] = "NIP PNS harus 18 digit angka"
	}

	minDate := time.Date(1940, 1, 1, 0, 0, 0, 0, time.UTC)
	switch {
	case r.TanggalLahir.IsZero():
		errs["tanggal_lahir"] = "tanggal lahir wajib diisi (YYYY-MM-DD)"
	case r.TanggalLahir.After(time.Now()):
		errs["tanggal_lahir"] = "tanggal lahir tidak boleh di masa depan"
	case r.TanggalLahir.Before(minDate):
		errs["tanggal_lahir"] = "tanggal lahir minimal 1940-01-01"
	}
	return errs
}

// ApplyTo: salin isi request ke model.
func (r *PersonelRequest) ApplyTo(p *model.Personel) {
	p.Jenis = r.Jenis
	p.NrpNip = r.NrpNip
	p.Nama = strings.TrimSpace(r.Nama)
	p.PangkatID = r.PangkatID
	p.SatkerID = r.SatkerID
	p.TempatLahir = r.TempatLahir
	p.TanggalLahir = r.TanggalLahir
	p.Status = r.Status
	if p.Status == "" {
		p.Status = "AKTIF"
	}
}
