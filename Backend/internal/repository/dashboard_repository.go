package repository

import (
	"time"

	"gorm.io/gorm"
)

type angka struct {
	Total        int64 `json:"total"`
	Polri        int64 `json:"polri"`
	Pns          int64 `json:"pns"`
	TanpaJabatan int64 `json:"tanpa_jabatan"`
	AkanPensiun  int64 `json:"akan_pensiun"`
}

type Jumlah struct {
	Label string `json:"label"`
	Total int64  `json:"total"`
}

type JumlahSatker struct {
	SatkerID    int    `json:"satker_id"`
	Nama        string `json:"nama"`
	Kode        string `json:"kode"`
	Total       int64  `json:"total"`
	Polri       int64  `json:"polri"`
	Pns         int64  `json:"pns"`
	AkanPensiun int64  `json:"akan_pensiun"`
}

type JumlahPangkat struct {
	Jenis    string `json:"jenis"`
	Kelompok string `json:"kelompok"`
	Total    int64  `json:"total"`
}

type PersonelSingkat struct {
	ID      int64  `json:"id"`
	Nama    string `json:"nama"`
	NrpNip  string `json:"nrp_nip"`
	Pangkat string `json:"pangkat"`
	Satker  string `json:"satker"`
}

type Pensiun struct {
	ID             int64     `json:"id"`
	Nama           string    `json:"nama"`
	Pangkat        string    `json:"pangkat"`
	Satker         string    `json:"satker"`
	Jabatan        *string   `json:"jabatan"`
	TanggalPensiun time.Time `json:"tanggal_pensiun"`
}

type Mutasi struct {
	PersonelID  int64     `json:"personel_id"`
	Nama        string    `json:"nama"`
	Pangkat     string    `json:"pangkat"`
	NamaJabatan string    `json:"nama_jabatan"`
	Satker      string    `json:"satker"`
	TmtMulai    time.Time `json:"tmt_mulai"`
}

// Ringkasan beranda. Field slice dipisah dari `angka` supaya GORM
// tidak menganggapnya relasi tabel saat Scan.
type Ringkasan struct {
	angka
	PerSatker    []JumlahSatker    `json:"per_satker"`
	PerPangkat   []JumlahPangkat   `json:"per_kelompok_pangkat"`
	PerluJabatan []PersonelSingkat `json:"perlu_jabatan"`
	AkanPensiun  []Pensiun         `json:"akan_pensiun_list"`
	Mutasi       []Mutasi          `json:"mutasi_terbaru"`
}

type DashboardRepository interface {
	Ringkasan(scopeSatkerID *int) (*Ringkasan, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// Batas usia pensiun 58 tahun; dihitung yang jatuh dalam 12 bulan ke depan.
const akanPensiun = `p.status = 'AKTIF' AND (p.tanggal_lahir + interval '58 years')::date
	BETWEEN CURRENT_DATE AND (CURRENT_DATE + interval '12 months')::date`

const tanpaJabatanAktif = `NOT EXISTS (
	SELECT 1 FROM riwayat_jabatan rj
	WHERE rj.personel_id = p.id AND rj.tmt_selesai IS NULL AND rj.deleted_at IS NULL)`

func (r *dashboardRepository) Ringkasan(scopeSatkerID *int) (*Ringkasan, error) {
	where := "p.deleted_at IS NULL"
	var args []any
	if scopeSatkerID != nil {
		where += " AND p.satker_id IN (SELECT id FROM satker_turunan(?))"
		args = append(args, *scopeSatkerID)
	}

	var res Ringkasan
	if err := r.db.Raw(`
		SELECT count(*) AS total,
		       count(*) FILTER (WHERE p.jenis = 'POLRI') AS polri,
		       count(*) FILTER (WHERE p.jenis = 'PNS')   AS pns,
		       count(*) FILTER (WHERE p.status = 'AKTIF' AND `+tanpaJabatanAktif+`) AS tanpa_jabatan,
		       count(*) FILTER (WHERE `+akanPensiun+`) AS akan_pensiun
		FROM personel p WHERE `+where, args...).Scan(&res.angka).Error; err != nil {
		return nil, err
	}

	// Sebaran per satker satu tingkat di bawah wilayah pemanggil
	// (operator PMJ → per Polres/Direktorat; admin → per Polda/Satker Mabes).
	var root any = gorm.Expr("(SELECT id FROM satker WHERE parent_id IS NULL ORDER BY id LIMIT 1)")
	if scopeSatkerID != nil {
		root = *scopeSatkerID
	}
	if err := r.db.Raw(`
		WITH RECURSIVE cabang AS (
			SELECT s.id, s.id AS induk FROM satker s WHERE s.parent_id = ?
			UNION ALL
			SELECT c.id, b.induk FROM satker c JOIN cabang b ON c.parent_id = b.id
		)
		SELECT s.id AS satker_id, s.nama, s.kode, count(*) AS total,
		       count(*) FILTER (WHERE p.jenis = 'POLRI') AS polri,
		       count(*) FILTER (WHERE p.jenis = 'PNS')   AS pns,
		       count(*) FILTER (WHERE `+akanPensiun+`) AS akan_pensiun
		FROM personel p
		JOIN cabang b ON b.id = p.satker_id
		JOIN satker s ON s.id = b.induk
		WHERE p.deleted_at IS NULL
		GROUP BY s.id, s.nama, s.kode
		ORDER BY total DESC`, root).Scan(&res.PerSatker).Error; err != nil {
		return nil, err
	}

	if err := r.db.Raw(`
		SELECT pg.jenis, pg.kelompok, count(*) AS total
		FROM personel p JOIN pangkat pg ON pg.id = p.pangkat_id
		WHERE `+where+`
		GROUP BY pg.jenis, pg.kelompok
		ORDER BY pg.jenis DESC, max(pg.urutan) DESC`, args...).Scan(&res.PerPangkat).Error; err != nil {
		return nil, err
	}

	if err := r.db.Raw(`
		SELECT p.id, p.nama, p.nrp_nip, pg.kode AS pangkat, s.nama AS satker
		FROM personel p
		JOIN pangkat pg ON pg.id = p.pangkat_id
		JOIN satker s ON s.id = p.satker_id
		WHERE `+where+` AND p.status = 'AKTIF' AND `+tanpaJabatanAktif+`
		ORDER BY p.nama LIMIT 6`, args...).Scan(&res.PerluJabatan).Error; err != nil {
		return nil, err
	}

	if err := r.db.Raw(`
		SELECT p.id, p.nama, pg.kode AS pangkat, s.nama AS satker,
		       (SELECT rj.nama_jabatan FROM riwayat_jabatan rj
		        WHERE rj.personel_id = p.id AND rj.tmt_selesai IS NULL AND rj.deleted_at IS NULL
		        ORDER BY rj.tmt_mulai DESC LIMIT 1) AS jabatan,
		       (p.tanggal_lahir + interval '58 years')::date AS tanggal_pensiun
		FROM personel p
		JOIN pangkat pg ON pg.id = p.pangkat_id
		JOIN satker s ON s.id = p.satker_id
		WHERE `+where+` AND `+akanPensiun+`
		ORDER BY tanggal_pensiun LIMIT 6`, args...).Scan(&res.AkanPensiun).Error; err != nil {
		return nil, err
	}

	if err := r.db.Raw(`
		SELECT p.id AS personel_id, p.nama, pg.kode AS pangkat, rj.nama_jabatan, s.nama AS satker, rj.tmt_mulai
		FROM riwayat_jabatan rj
		JOIN personel p ON p.id = rj.personel_id
		JOIN pangkat pg ON pg.id = p.pangkat_id
		JOIN satker s ON s.id = rj.satker_id
		WHERE `+where+` AND rj.deleted_at IS NULL AND rj.tmt_mulai <= CURRENT_DATE
		ORDER BY rj.tmt_mulai DESC, rj.id DESC LIMIT 6`, args...).Scan(&res.Mutasi).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
