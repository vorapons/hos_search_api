package repository

import (
	"errors"
	"pt_search_hos/domain"

	"gorm.io/gorm"
)

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) domain.StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) FindByEmail(email string) (*domain.Staff, error) {
	var m StaffModel
	err := r.db.Preload("Hospital").Where("email = ?", email).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &domain.Staff{
		ID:           m.ID,
		HospitalID:   m.HospitalID,
		HospitalName: m.Hospital.Name,
		Email:        m.Email,
		Password:     m.Password,
	}, nil
}

func (r *staffRepository) Create(staff *domain.Staff) error {
	m := StaffModel{
		HospitalID: staff.HospitalID,
		Email:      staff.Email,
		Password:   staff.Password,
	}
	return r.db.Create(&m).Error
}

func (r *staffRepository) FindHospitalByName(name string) (*domain.Hospital, error) {
	var m HospitalModel
	err := r.db.Where("name = ?", name).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &domain.Hospital{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}
