package service

import (
	"Backend/internal/dto"
	"Backend/internal/model"
	"Backend/internal/repository"
	"errors"
)

// Scope riwayat jabatan mengikuti satker PERSONEL-nya (bukan satker di jabatan),
// karena riwayat bisa berisi jabatan di satker lain sebelum mutasi.
type RiwayatJabatanService struct {
	repo         repository.RiwayatJabatanRepository
	personelRepo repository.PersonelRepository
}

func NewRiwayatJabatanService(repo repository.RiwayatJabatanRepository, personelRepo repository.PersonelRepository) *RiwayatJabatanService {
	return &RiwayatJabatanService{repo: repo, personelRepo: personelRepo}
}

func (s *RiwayatJabatanService) List(p *model.Principal, personelID int64) ([]model.RiwayatJabatan, error) {
	if _, err := s.personelRepo.FindBasic(personelID, p.SatkerID); err != nil {
		return nil, err
	}
	return s.repo.ListByPersonel(personelID)
}

func (s *RiwayatJabatanService) Create(p *model.Principal, personelID int64, req dto.RiwayatJabatanRequest) (*model.RiwayatJabatan, error) {
	personel, err := s.personelRepo.FindBasic(personelID, p.SatkerID)
	if err != nil {
		return nil, err
	}
	if errs := req.Validate(personel.TanggalLahir); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}

	rj := model.RiwayatJabatan{PersonelID: personelID}
	req.ApplyTo(&rj)
	if err := s.repo.Create(&rj); err != nil {
		return nil, err
	}
	return s.repo.FindByID(rj.ID)
}

func (s *RiwayatJabatanService) Update(p *model.Principal, id int64, req dto.RiwayatJabatanRequest) (*model.RiwayatJabatan, error) {
	rj, personel, err := s.authorize(p, id)
	if err != nil {
		return nil, err
	}
	if errs := req.Validate(personel.TanggalLahir); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}

	req.ApplyTo(rj)
	rj.Satker, rj.Fungsi, rj.Nivelering = nil, nil, nil // jangan ikut tersimpan
	if err := s.repo.Update(rj); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

func (s *RiwayatJabatanService) Delete(p *model.Principal, id int64) error {
	if _, _, err := s.authorize(p, id); err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}

// authorize: riwayat harus ada DAN personel pemiliknya dalam scope user.
// Di luar scope dijawab "riwayat tidak ditemukan" (404), bukan 403.
func (s *RiwayatJabatanService) authorize(p *model.Principal, id int64) (*model.RiwayatJabatan, *model.Personel, error) {
	rj, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	personel, err := s.personelRepo.FindBasic(rj.PersonelID, p.SatkerID)
	if errors.Is(err, repository.ErrPersonelNotFound) {
		return nil, nil, repository.ErrRiwayatNotFound
	}
	if err != nil {
		return nil, nil, err
	}
	return rj, personel, nil
}
