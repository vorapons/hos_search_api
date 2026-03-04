package domain

import (
	"time"
)

type Patient struct {
	ID           string     `json:"-"`
	HospitalID   string     `json:"-"`
	FirstNameTH  *string    `json:"first_name_th"`
	MiddleNameTH *string    `json:"middle_name_th"`
	LastNameTH   *string    `json:"last_name_th"`
	FirstNameEN  *string    `json:"first_name_en"`
	MiddleNameEN *string    `json:"middle_name_en"`
	LastNameEN   *string    `json:"last_name_en"`
	NationalID   *string    `json:"national_id"`
	PassportID   *string    `json:"passport_id"`
	PatientHN    *string    `json:"patient_hn"`
	DateOfBirth  *time.Time `json:"date_of_birth"`
	Gender       *string    `json:"gender"`
	PhoneNumber  *string    `json:"phone_number"`
	Email        *string    `json:"email"`
}

// PatientSearchInput holds search criteria from the request body.
type PatientSearchInput struct {
	NationalID  *string    `json:"national_id"`
	PassportID  *string    `json:"passport_id"`
	FirstName   *string    `json:"first_name"`
	MiddleName  *string    `json:"middle_name"`
	LastName    *string    `json:"last_name"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	PhoneNumber *string    `json:"phone_number"`
	Email       *string    `json:"email"`
}

// PatientRepository is the port (interface) for the database adapter.
type PatientRepository interface {
	FindByID(id string, hospitalID string) (*Patient, error)
	FindByCondition(input PatientSearchInput, hospitalID string) ([]Patient, error)
}

// PatientService is the port (interface) for the use-case layer.
type PatientService interface {
	GetPatientByID(id string, hospitalID string) (*Patient, error)
	GetPatientByCondition(input PatientSearchInput, hospitalID string) ([]Patient, error)
}
