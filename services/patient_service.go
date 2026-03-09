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

// allowedOrderBy mirrors the repo whitelist for validation.
var allowedOrderBy = map[string]bool{
	"first_name_th": true,
	"last_name_th":  true,
	"first_name_en": true,
	"last_name_en":  true,
	"date_of_birth": true,
	"patient_hn":    true,
}

func (s *patientService) GetPatientByCondition(input domain.PatientSearchInput, hospitalID string) (domain.PatientSearchResult, error) {
	if hospitalID == "" {
		return domain.PatientSearchResult{}, domain.ErrInvalidInput
	}
	if input.NationalID == nil && input.PassportID == nil &&
		input.FirstName == nil && input.MiddleName == nil && input.LastName == nil &&
		input.DateOfBirth == nil && input.PhoneNumber == nil && input.Email == nil {
		return domain.PatientSearchResult{}, domain.ErrInvalidInput
	}

	// Normalise pagination defaults
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	} else if input.PageSize > 100 {
		input.PageSize = 100
	}
	if !allowedOrderBy[input.OrderBy] {
		input.OrderBy = "last_name_th"
	}
	if input.OrderDir != "asc" && input.OrderDir != "desc" {
		input.OrderDir = "asc"
	}

	patients, total, err := s.repo.FindByCondition(input, hospitalID)
	if err != nil {
		return domain.PatientSearchResult{}, err
	}

	return domain.PatientSearchResult{
		Data:     patients,
		Total:    total,
		Page:     input.Page,
		PageSize: input.PageSize,
	}, nil
}
