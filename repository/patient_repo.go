package repository

import (
	"errors"

	"pt_search_hos/domain"

	"gorm.io/gorm"
)

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) domain.PatientRepository {
	return &patientRepository{db: db}
}

// FindByID searches for a patient by national_id OR passport_id within the given hospital.
func (r *patientRepository) FindByID(id string, hospitalID string) (*domain.Patient, error) {
	var m PatientModel
	err := r.db.
		Where("(national_id = ? OR passport_id = ?) AND hospital_id = ?", id, id, hospitalID).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toPatient(&m), nil
}

// allowedOrderBy maps safe client-facing field names to actual DB columns.
var allowedOrderBy = map[string]string{
	"first_name_th": "first_name_th",
	"last_name_th":  "last_name_th",
	"first_name_en": "first_name_en",
	"last_name_en":  "last_name_en",
	"date_of_birth": "date_of_birth",
	"patient_hn":    "patient_hn",
}

// FindByCondition returns a paginated, ordered page of patients in the given hospital
// that match all provided fields. All text fields use case-insensitive contains matching.
// DateOfBirth is an exact match.
// input.Page, input.PageSize, input.OrderBy, input.OrderDir must already be normalised by the service.
func (r *patientRepository) FindByCondition(input domain.PatientSearchInput, hospitalID string) ([]domain.Patient, int64, error) {
	q := r.db.Model(&PatientModel{}).Where("hospital_id = ?", hospitalID)

	if input.NationalID != nil {
		q = q.Where("national_id ILIKE ?", "%"+*input.NationalID+"%")
	}
	if input.PassportID != nil {
		q = q.Where("passport_id ILIKE ?", "%"+*input.PassportID+"%")
	}
	if input.FirstName != nil {
		like := "%" + *input.FirstName + "%"
		q = q.Where("(first_name_th ILIKE ? OR first_name_en ILIKE ?)", like, like)
	}
	if input.MiddleName != nil {
		like := "%" + *input.MiddleName + "%"
		q = q.Where("(middle_name_th ILIKE ? OR middle_name_en ILIKE ?)", like, like)
	}
	if input.LastName != nil {
		like := "%" + *input.LastName + "%"
		q = q.Where("(last_name_th ILIKE ? OR last_name_en ILIKE ?)", like, like)
	}
	if input.DateOfBirth != nil {
		q = q.Where("date_of_birth = ?", *input.DateOfBirth)
	}
	if input.PhoneNumber != nil {
		q = q.Where("phone_number ILIKE ?", "%"+*input.PhoneNumber+"%")
	}
	if input.Email != nil {
		q = q.Where("email ILIKE ?", "%"+*input.Email+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	col, ok := allowedOrderBy[input.OrderBy]
	if !ok {
		col = "last_name_th"
	}
	dir := "ASC"
	if input.OrderDir == "desc" {
		dir = "DESC"
	}

	offset := (input.Page - 1) * input.PageSize

	var models []PatientModel
	if err := q.Order(col + " " + dir).Limit(input.PageSize).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	patients := make([]domain.Patient, len(models))
	for i, m := range models {
		patients[i] = *toPatient(&m)
	}
	return patients, total, nil
}

func toPatient(m *PatientModel) *domain.Patient {
	return &domain.Patient{
		ID:           m.ID,
		HospitalID:   m.HospitalID,
		FirstNameTH:  m.FirstNameTH,
		MiddleNameTH: m.MiddleNameTH,
		LastNameTH:   m.LastNameTH,
		FirstNameEN:  m.FirstNameEN,
		MiddleNameEN: m.MiddleNameEN,
		LastNameEN:   m.LastNameEN,
		NationalID:   m.NationalID,
		PassportID:   m.PassportID,
		PatientHN:    m.PatientHN,
		DateOfBirth:  m.DateOfBirth,
		Gender:       m.Gender,
		PhoneNumber:  m.PhoneNumber,
		Email:        m.Email,
	}
}
