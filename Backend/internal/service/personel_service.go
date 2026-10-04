package service

import (
	"Backend/internal/dto"
	"Backend/internal/model"
	"Backend/internal/repository"
	"errors"
)

var ErrSatkerOutOfScope = errors.New("satker di luar kewenangan")

// ValidationError: kumpulan error per field, dikirim ke client sebagai 422.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string { return "validasi gagal" }

type PersonelService struct {
	repo repository.PersonelRepository
}

func NewPersonelService(repo repository.PersonelRepository) *PersonelService {
	return &PersonelService{repo: repo}
}

// Hasil list + info paginasi
type PersonelListResult struct {
	Items []model.Personel `json:"items"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Total int64            `json:"total"`
}

func (s *PersonelService) List(p *model.Principal, f repository.PersonelFilter) (*PersonelListResult, error) {
	items, total, err := s.repo.List(f, p.SatkerID)
	if err != nil {
		return nil, err
	}

	for i := range items {
		setJabatanSaatIni(&items[i])
		items[i].RiwayatJabatan = nil // list gak perlu riwayat, cukup jabatan_saat_ini
	}

	// samakan dengan normalisasi di repo
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	return &PersonelListResult{Items: items, Page: f.Page, Limit: f.Limit, Total: total}, nil
}

func (s *PersonelService) GetByID(p *model.Principal, id int64) (*model.Personel, error) {
	personel, err := s.repo.FindByID(id, p.SatkerID)
	if err != nil {
		return nil, err
	}
	setJabatanSaatIni(personel)
	return personel, nil
}

// Jabatan saat ini = DEFINITIF yang tmt_selesai-nya kosong.
// Kalau gak ada (mis. cuma PLT/PLH), ambil jabatan aktif terakhir.
func setJabatanSaatIni(p *model.Personel) {
	var fallback *model.RiwayatJabatan
	for i := range p.RiwayatJabatan {
		rj := &p.RiwayatJabatan[i]
		if rj.TmtSelesai != nil {
			continue
		}
		if rj.StatusJabatan == "DEFINITIF" {
			p.JabatanSaatIni = rj
			return
		}
		fallback = rj
	}
	p.JabatanSaatIni = fallback
}

// Create: validasi → cek satker tujuan dalam scope → simpan.
func (s *PersonelService) Create(p *model.Principal, req dto.PersonelRequest) (*model.Personel, error) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}
	if err := s.checkSatker(p, req.SatkerID); err != nil {
		return nil, err
	}

	var personel model.Personel
	req.ApplyTo(&personel)
	if err := s.repo.Create(&personel); err != nil {
		return nil, err
	}
	return s.GetByID(p, personel.ID)
}

// Update: validasi → personel lama harus dalam scope → satker baru juga harus dalam scope.
func (s *PersonelService) Update(p *model.Principal, id int64, req dto.PersonelRequest) (*model.Personel, error) {
	if errs := req.Validate(); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}

	existing, err := s.repo.FindBasic(id, p.SatkerID) // scope cek #1
	if err != nil {
		return nil, err
	}
	if err := s.checkSatker(p, req.SatkerID); err != nil { // scope cek #2
		return nil, err
	}

	req.ApplyTo(existing)
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return s.GetByID(p, id)
}

// Delete: personel harus dalam scope → soft delete.
func (s *PersonelService) Delete(p *model.Principal, id int64) error {
	if _, err := s.repo.FindBasic(id, p.SatkerID); err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}

func (s *PersonelService) checkSatker(p *model.Principal, satkerID int) error {
	ok, err := s.repo.SatkerInScope(satkerID, p.SatkerID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrSatkerOutOfScope
	}
	return nil
}
