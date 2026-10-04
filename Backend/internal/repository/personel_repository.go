package repository

import (
	"Backend/internal/model"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrPersonelNotFound = errors.New("personel tidak ditemukan")

// Parameter pencarian dari query string
type PersonelFilter struct {
	Q            string // cari nama / NRP
	SatkerID     *int   // filter satker (ikut turunannya)
	Status       string // AKTIF, PENSIUN, ...
	Jenis        string // POLRI / PNS
	TanpaJabatan bool   // hanya yang belum punya jabatan aktif
	Kelompok     string // kelompok pangkat, mis. "Bintara", "Golongan III"
	AkanPensiun  bool   // aktif & mencapai usia 58 tahun dalam 12 bulan ke depan
	Urut         string // nama, -nama, pangkat, -pangkat, satker, -satker, pensiun, -pensiun
	Page         int
	Limit        int
}

type PersonelRepository interface {
	List(f PersonelFilter, scopeSatkerID *int) ([]model.Personel, int64, error)
	FindByID(id int64, scopeSatkerID *int) (*model.Personel, error)
	FindBasic(id int64, scopeSatkerID *int) (*model.Personel, error)
	SatkerInScope(satkerID int, scopeSatkerID *int) (bool, error)
	Create(p *model.Personel) error
	Update(p *model.Personel) error
	SoftDelete(id int64) error
}

type personelRepository struct {
	db *gorm.DB
}

func NewPersonelRepository(db *gorm.DB) PersonelRepository {
	return &personelRepository{db: db}
}

// scope: operator cuma boleh lihat satker-nya + turunannya.
// scopeSatkerID nil = admin, tanpa batasan.
func applyScope(q *gorm.DB, scopeSatkerID *int) *gorm.DB {
	if scopeSatkerID != nil {
		q = q.Where("personel.satker_id IN (SELECT id FROM satker_turunan(?))", *scopeSatkerID)
	}
	return q
}

func (r *personelRepository) List(f PersonelFilter, scopeSatkerID *int) ([]model.Personel, int64, error) {
	// normalisasi paginasi
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	q := r.db.Model(&model.Personel{})
	q = applyScope(q, scopeSatkerID)

	if s := strings.TrimSpace(f.Q); s != "" {
		q = q.Where("lower(personel.nama) LIKE ? OR personel.nrp_nip LIKE ?",
			"%"+strings.ToLower(s)+"%", s+"%")
	}
	if f.SatkerID != nil {
		q = q.Where("personel.satker_id IN (SELECT id FROM satker_turunan(?))", *f.SatkerID)
	}
	if f.Status != "" {
		q = q.Where("personel.status = ?", f.Status)
	}
	if f.Jenis != "" {
		q = q.Where("personel.jenis = ?", f.Jenis)
	}
	if f.Kelompok != "" {
		q = q.Where("personel.pangkat_id IN (SELECT id FROM pangkat WHERE kelompok = ?)", f.Kelompok)
	}
	if f.TanpaJabatan {
		q = q.Where(`NOT EXISTS (SELECT 1 FROM riwayat_jabatan rj
			WHERE rj.personel_id = personel.id AND rj.tmt_selesai IS NULL AND rj.deleted_at IS NULL)`)
	}

	if f.AkanPensiun {
		q = q.Where("personel.status = 'AKTIF' AND " + batasPensiun("personel"))
	}
	urutan := urutanPersonel(f, "personel")

	// hitung total dulu (sebelum limit/offset)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Personel
	err := q.
		Preload("Pangkat").
		Preload("Satker").
		// list cuma butuh jabatan aktif, bukan seluruh riwayat
		Preload("RiwayatJabatan", "tmt_selesai IS NULL AND status_jabatan = ?", "DEFINITIF").
		Order(urutan).
		Limit(f.Limit).
		Offset((f.Page - 1) * f.Limit).
		Find(&items).Error

	return items, total, err
}

func (r *personelRepository) FindByID(id int64, scopeSatkerID *int) (*model.Personel, error) {
	q := r.db.
		Preload("Pangkat").
		Preload("Satker").
		Preload("RiwayatJabatan", func(db *gorm.DB) *gorm.DB {
			return db.Order("tmt_mulai ASC") // kronologis: lama → baru
		}).
		Preload("RiwayatJabatan.Satker").
		Preload("RiwayatJabatan.Fungsi").
		Preload("RiwayatJabatan.Nivelering").
		// sertifikasi = kualifikasi jenis PELATIHAN, terbaru dulu
		Preload("Sertifikasi", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN master_pendidikan mp ON mp.id = kualifikasi.master_pendidikan_id").
				Where("mp.jenis = ?", "PELATIHAN").
				Order("kualifikasi.tahun_lulus DESC")
		}).
		Preload("Sertifikasi.Pendidikan")

	q = applyScope(q, scopeSatkerID)

	var p model.Personel
	err := q.First(&p, "personel.id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPersonelNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBasic: tanpa preload — cukup untuk cek "ada & masih dalam scope".
func (r *personelRepository) FindBasic(id int64, scopeSatkerID *int) (*model.Personel, error) {
	var p model.Personel
	err := applyScope(r.db, scopeSatkerID).First(&p, "personel.id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPersonelNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// SatkerInScope: apakah satkerID termasuk wilayah kewenangan user.
func (r *personelRepository) SatkerInScope(satkerID int, scopeSatkerID *int) (bool, error) {
	if scopeSatkerID == nil {
		return true, nil // admin: semua satker
	}
	var ok bool
	err := r.db.Raw(
		"SELECT EXISTS (SELECT 1 FROM satker_turunan(?) WHERE id = ?)",
		*scopeSatkerID, satkerID,
	).Scan(&ok).Error
	return ok, err
}

func (r *personelRepository) Create(p *model.Personel) error {
	return r.db.Omit(clause.Associations).Create(p).Error
}

// Update: Select kolom eksplisit → nilai kosong tetap ter-update,
// dan relasi (Pangkat, Satker, RiwayatJabatan) tidak ikut disentuh.
func (r *personelRepository) Update(p *model.Personel) error {
	return r.db.Model(&model.Personel{ID: p.ID}).
		Select("jenis", "nrp_nip", "nama", "pangkat_id", "satker_id",
			"tempat_lahir", "tanggal_lahir", "status").
		Updates(p).Error
}

// SoftDelete: isi deleted_at personel + seluruh riwayat jabatannya dalam 1 transaksi.
func (r *personelRepository) SoftDelete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("personel_id = ?", id).Delete(&model.RiwayatJabatan{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Personel{}, id).Error
	})
}

// batasPensiun: usia 58 tahun jatuh antara hari ini dan 12 bulan ke depan.
func batasPensiun(tabel string) string {
	return "(" + tabel + ".tanggal_lahir + interval '58 years')::date " +
		"BETWEEN CURRENT_DATE AND (CURRENT_DATE + interval '12 months')::date"
}

// urutanPersonel: ORDER BY dari parameter sort (whitelist, aman dari injeksi).
// Default: nama A-Z, atau yang paling dekat pensiun bila filter akan pensiun aktif.
func urutanPersonel(f PersonelFilter, t string) string {
	urut := f.Urut
	if urut == "" && f.AkanPensiun {
		urut = "pensiun"
	}
	arah := "ASC"
	kunci := urut
	if strings.HasPrefix(urut, "-") {
		arah, kunci = "DESC", urut[1:]
	}
	balik := map[string]string{"ASC": "DESC", "DESC": "ASC"}
	switch kunci {
	case "pangkat": // pangkat tertinggi dulu
		return "(SELECT pg.urutan FROM pangkat pg WHERE pg.id = " + t + ".pangkat_id) " + balik[arah] + ", " + t + ".nama"
	case "satker":
		return "(SELECT lower(s.nama) FROM satker s WHERE s.id = " + t + ".satker_id) " + arah + ", " + t + ".nama"
	case "pensiun": // paling tua = paling dekat pensiun
		return t + ".tanggal_lahir " + arah + ", " + t + ".nama"
	case "nama":
		return t + ".nama " + arah
	}
	return t + ".nama ASC"
}
