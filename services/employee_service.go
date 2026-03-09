package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type EmployeeService struct {
	repo *repositories.EmployeeRepository
}

func NewEmployeeService(repo *repositories.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) GetAll(aktif *bool) ([]models.Employee, error) {
	return s.repo.GetAll(aktif)
}

func (s *EmployeeService) GetByID(id int) (*models.Employee, error) {
	employee, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}
	return employee, nil
}

func (s *EmployeeService) Create(req models.CreateEmployeeRequest) (*models.Employee, error) {
	// Parse tanggal_masuk string to time.Time
	var tanggalMasuk *time.Time
	if req.TanggalMasuk != nil && *req.TanggalMasuk != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalMasuk)
		if err != nil {
			return nil, errors.New("format tanggal_masuk tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalMasuk = &t
	}

	emp := &models.Employee{
		Nama:         req.Nama,
		Posisi:       req.Posisi,
		GajiPokok:    req.GajiPokok,
		NoHp:         req.NoHp,
		Alamat:       req.Alamat,
		TanggalMasuk: tanggalMasuk,
		UserID:       req.UserID,
	}

	err := s.repo.Create(emp)
	if err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *EmployeeService) Update(id int, req models.UpdateEmployeeRequest) (*models.Employee, error) {
	emp, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}

	if req.Nama != "" {
		emp.Nama = req.Nama
	}
	if req.Posisi != "" {
		emp.Posisi = req.Posisi
	}
	if req.GajiPokok != nil {
		emp.GajiPokok = *req.GajiPokok
	}
	if req.NoHp != nil {
		emp.NoHp = req.NoHp
	}
	if req.Alamat != nil {
		emp.Alamat = req.Alamat
	}
	if req.UserID != nil {
		emp.UserID = req.UserID
	}
	if req.Aktif != nil {
		emp.Aktif = *req.Aktif
	}
	if req.TanggalMasuk != nil && *req.TanggalMasuk != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalMasuk)
		if err == nil {
			emp.TanggalMasuk = &t
		}
	}

	err = s.repo.Update(emp)
	if err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *EmployeeService) SoftDelete(id int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("karyawan tidak ditemukan")
	}
	return s.repo.SoftDelete(id)
}
