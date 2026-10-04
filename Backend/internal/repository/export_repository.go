package repository

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type BarisPersonel struct {
	NrpNip       string     `json:"nrp_nip"`
	Nama         string     `json:"nama"`
	Jenis        string     `json:"jenis"`
	Pangkat      string     `json:"pangkat"`
	Satker       string     `json:"satker"`
	Jabatan      *string    `json:"jabatan"`
	TmtJabatan   *time.Time `json:"tmt_jabatan"`
	TempatLahir  *string    `json:"tempat_lahir"`
	TanggalLahir time.Time  `json:"tanggal_lahir"`
	Status       string     `json:"status"`
}

type ExportRepository interface {
	DaftarPersonel(f PersonelFilter, scopeSatkerID *int) ([]BarisPersonel, error)
}

type exportRepository struct {
	db *gorm.DB
}

func NewExportRepository(db *gorm.DB) ExportRepository {
	return &exportRepository{db: db}
}

// DaftarPersonel: personel di wilayah pemanggil + jabatan saat ini,
// dengan filter yang sama seperti halaman Data Personel.
func (r *exportRepository) DaftarPersonel(f PersonelFilter, scopeSatkerID *int) ([]BarisPersonel, error) {
	where := "p.deleted_at IS NULL"
	var args []any
	if scopeSatkerID != nil {
		where += " AND p.satker_id IN (SELECT id FROM satker_turunan(?))"
		args = append(args, *scopeSatkerID)
	}
	if s := strings.TrimSpace(f.Q); s != "" {
		where += " AND (lower(p.nama) LIKE ? OR p.nrp_nip LIKE ?)"
		args = append(args, "%"+strings.ToLower(s)+"%", s+"%")
	}
	if f.SatkerID != nil {
		where += " AND p.satker_id IN (SELECT id FROM satker_turunan(?))"
		args = append(args, *f.SatkerID)
	}
	if f.Status != "" {
		where += " AND p.status = ?"
		args = append(args, f.Status)
	}
	if f.Jenis != "" {
		where += " AND p.jenis = ?"
		args = append(args, f.Jenis)
	}
	if f.Kelompok != "" {
		where += " AND pg.kelompok = ?"
		args = append(args, f.Kelompok)
	}
	if f.TanpaJabatan {
		where += " AND j.nama_jabatan IS NULL"
	}
	urutan := "s.nama, pg.jenis DESC, pg.urutan DESC, p.nama"
	if f.AkanPensiun {
		where += " AND p.status = 'AKTIF' AND " + batasPensiun("p")
	}
	if f.Urut != "" || f.AkanPensiun {
		urutan = urutanPersonel(f, "p") // sama dengan urutan di layar
	}
	var rows []BarisPersonel
	err := r.db.Raw(`
		SELECT p.nrp_nip, p.nama, p.jenis, pg.kode AS pangkat, s.nama AS satker,
		       j.nama_jabatan AS jabatan, j.tmt_mulai AS tmt_jabatan,
		       p.tempat_lahir, p.tanggal_lahir, p.status
		FROM personel p
		JOIN pangkat pg ON pg.id = p.pangkat_id
		JOIN satker s ON s.id = p.satker_id
		LEFT JOIN LATERAL (
			SELECT rj.nama_jabatan, rj.tmt_mulai FROM riwayat_jabatan rj
			WHERE rj.personel_id = p.id AND rj.tmt_selesai IS NULL AND rj.deleted_at IS NULL
			ORDER BY (rj.status_jabatan = 'DEFINITIF') DESC, rj.tmt_mulai DESC LIMIT 1
		) j ON true
		WHERE `+where+`
		ORDER BY `+urutan, args...).Scan(&rows).Error
	return rows, err
}
