package services_test

import (
	"testing"
	"time"

	"pt_search_hos/domain"
	"pt_search_hos/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ── mock repo ─────────────────────────────────────────────────────────────────

type mockPatientRepo struct {
	mock.Mock
}

func (m *mockPatientRepo) FindByID(id, hospitalID string) (*domain.Patient, error) {
	args := m.Called(id, hospitalID)
	patient, _ := args.Get(0).(*domain.Patient)
	return patient, args.Error(1)
}

func (m *mockPatientRepo) FindByCondition(input domain.PatientSearchInput, hospitalID string) ([]domain.Patient, error) {
	args := m.Called(input, hospitalID)
	patients, _ := args.Get(0).([]domain.Patient)
	return patients, args.Error(1)
}

// ── helper ────────────────────────────────────────────────────────────────────

func ptr(s string) *string { return &s }

func newPatientService(repo domain.PatientRepository) domain.PatientService {
	return services.NewPatientService(repo)
}

// ── GetPatientByID ────────────────────────────────────────────────────────────

// positive: patient found by national ID → returns patient
func TestGetPatientByID_FoundByNationalID(t *testing.T) {
	repo := new(mockPatientRepo)
	patient := &domain.Patient{ID: "uuid-1", NationalID: ptr("1234567890123")}
	repo.On("FindByID", "1234567890123", "BKH01").Return(patient, nil)

	result, err := newPatientService(repo).GetPatientByID("1234567890123", "BKH01")

	assert.NoError(t, err)
	assert.Equal(t, "uuid-1", result.ID)
	repo.AssertExpectations(t)
}

// positive: patient found by passport ID → returns patient
func TestGetPatientByID_FoundByPassportID(t *testing.T) {
	repo := new(mockPatientRepo)
	patient := &domain.Patient{ID: "uuid-2", PassportID: ptr("AB123456")}
	repo.On("FindByID", "AB123456", "BKH01").Return(patient, nil)

	result, err := newPatientService(repo).GetPatientByID("AB123456", "BKH01")

	assert.NoError(t, err)
	assert.Equal(t, "uuid-2", result.ID)
}

// negative: patient not in DB → ErrNotFound
func TestGetPatientByID_NotFound(t *testing.T) {
	repo := new(mockPatientRepo)
	repo.On("FindByID", "unknown", "BKH01").Return(nil, nil)

	_, err := newPatientService(repo).GetPatientByID("unknown", "BKH01")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

// negative: empty patient ID → ErrInvalidInput
func TestGetPatientByID_EmptyID(t *testing.T) {
	_, err := newPatientService(new(mockPatientRepo)).GetPatientByID("", "BKH01")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// negative: empty hospital ID → ErrInvalidInput
func TestGetPatientByID_EmptyHospitalID(t *testing.T) {
	_, err := newPatientService(new(mockPatientRepo)).GetPatientByID("1234567890123", "")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// negative: repository returns DB error → propagated as-is
func TestGetPatientByID_DBError(t *testing.T) {
	repo := new(mockPatientRepo)
	repo.On("FindByID", "1234567890123", "BKH01").Return(nil, assert.AnError)

	_, err := newPatientService(repo).GetPatientByID("1234567890123", "BKH01")
	assert.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrNotFound)
}

// ── GetPatientByCondition ─────────────────────────────────────────────────────

// positive: valid condition, results returned
func TestGetPatientByCondition_Success(t *testing.T) {
	repo := new(mockPatientRepo)
	input := domain.PatientSearchInput{LastName: ptr("Smith")}
	patients := []domain.Patient{{ID: "uuid-1"}, {ID: "uuid-2"}}
	repo.On("FindByCondition", input, "BKH01").Return(patients, nil)

	result, err := newPatientService(repo).GetPatientByCondition(input, "BKH01")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

// positive: valid condition, no patients found → empty slice
func TestGetPatientByCondition_EmptyResult(t *testing.T) {
	repo := new(mockPatientRepo)
	input := domain.PatientSearchInput{FirstName: ptr("NoOne")}
	repo.On("FindByCondition", input, "BKH01").Return([]domain.Patient{}, nil)

	result, err := newPatientService(repo).GetPatientByCondition(input, "BKH01")

	assert.NoError(t, err)
	assert.Empty(t, result)
}

// negative: all fields nil → ErrInvalidInput
func TestGetPatientByCondition_NoCondition(t *testing.T) {
	_, err := newPatientService(new(mockPatientRepo)).
		GetPatientByCondition(domain.PatientSearchInput{}, "BKH01")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// negative: empty hospital ID → ErrInvalidInput
func TestGetPatientByCondition_EmptyHospitalID(t *testing.T) {
	input := domain.PatientSearchInput{NationalID: ptr("1234567890123")}
	_, err := newPatientService(new(mockPatientRepo)).GetPatientByCondition(input, "")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// positive: search by national ID → results returned
func TestGetPatientByCondition_ByNationalID(t *testing.T) {
	repo := new(mockPatientRepo)
	input := domain.PatientSearchInput{NationalID: ptr("1234567890123")}
	repo.On("FindByCondition", input, "BKH01").Return([]domain.Patient{{ID: "uuid-1"}}, nil)

	result, err := newPatientService(repo).GetPatientByCondition(input, "BKH01")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// positive: search by date of birth → results returned
func TestGetPatientByCondition_ByDateOfBirth(t *testing.T) {
	repo := new(mockPatientRepo)
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	input := domain.PatientSearchInput{DateOfBirth: &dob}
	repo.On("FindByCondition", input, "BKH01").Return([]domain.Patient{{ID: "uuid-3"}}, nil)

	result, err := newPatientService(repo).GetPatientByCondition(input, "BKH01")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

// negative: repository returns DB error → propagated as-is
func TestGetPatientByCondition_DBError(t *testing.T) {
	repo := new(mockPatientRepo)
	input := domain.PatientSearchInput{LastName: ptr("Smith")}
	repo.On("FindByCondition", input, "BKH01").Return(nil, assert.AnError)

	_, err := newPatientService(repo).GetPatientByCondition(input, "BKH01")
	assert.Error(t, err)
}
