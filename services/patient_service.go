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
	if id == "" || hospitalID == "" {
		return nil, domain.ErrInvalidInput
	}

	patient, err := s.repo.FindByID(id, hospitalID)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, domain.ErrNotFound
	}
	return patient, nil
}

func (s *patientService) GetPatientByCondition(input domain.PatientSearchInput, hospitalID string) ([]domain.Patient, error) {
	if hospitalID == "" {
		return nil, domain.ErrInvalidInput
	}
	if input.NationalID == nil && input.PassportID == nil &&
		input.FirstName == nil && input.MiddleName == nil && input.LastName == nil &&
		input.DateOfBirth == nil && input.PhoneNumber == nil && input.Email == nil {
		return nil, domain.ErrInvalidInput
	}

	return s.repo.FindByCondition(input, hospitalID)
}
