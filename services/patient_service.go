package services

import (
	"pt_search_hos/domain"
)

type patientService struct {
	repo domain.PatientRepository
}

func NewPatientService(repo domain.PatientRepository) domain.PatientService {
	return &patientService{repo: repo}
}

func (s *patientService) GetPatientByID(id string, hospitalID string) (*domain.Patient, error) {
	// TODO: implement
	return nil, nil
}

func (s *patientService) GetPatientByCondition(input domain.PatientSearchInput, hospitalID string) ([]domain.Patient, error) {
	// TODO: implement
	return nil, nil
}
