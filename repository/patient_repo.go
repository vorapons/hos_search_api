package repository

import (
	"pt_search_hos/domain"

	"gorm.io/gorm"
)

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) domain.PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) FindByID(id string, hospitalID string) (*domain.Patient, error) {
	// TODO: implement
	return nil, nil
}

func (r *patientRepository) FindByCondition(input domain.PatientSearchInput, hospitalID string) ([]domain.Patient, error) {
	// TODO: implement
	return nil, nil
}
