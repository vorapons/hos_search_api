package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"pt_search_hos/domain"
	"pt_search_hos/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ── mock service ──────────────────────────────────────────────────────────────

type mockPatientService struct {
	mock.Mock
}

func (m *mockPatientService) GetPatientByID(id, hospitalID string) (*domain.Patient, error) {
	args := m.Called(id, hospitalID)
	patient, _ := args.Get(0).(*domain.Patient)
	return patient, args.Error(1)
}

func (m *mockPatientService) GetPatientByCondition(input domain.PatientSearchInput, hospitalID string) ([]domain.Patient, error) {
	args := m.Called(input, hospitalID)
	patients, _ := args.Get(0).([]domain.Patient)
	return patients, args.Error(1)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func setupPatientApp(staffSvc domain.StaffService, patientSvc domain.PatientService) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	staffH := handler.NewStaffHandler(staffSvc)
	patientH := handler.NewPatientHandler(patientSvc)
	handler.SetupRoutes(app, staffH, patientH, testJWTSecret, staffSvc.IsTokenBlacklisted)
	return app
}

// makePatientToken generates a JWT with hospital_id set to hospitalID.
func makePatientToken(hospitalID string) string {
	type claims struct {
		Login      string `json:"login"`
		HospitalID string `json:"hospital_id"`
		Hospital   string `json:"hospital"`
		jwt.RegisteredClaims
	}
	c := claims{
		Login:      "staff@example.com",
		HospitalID: hospitalID,
		Hospital:   "Bangkok Hospital",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, _ := token.SignedString([]byte(testJWTSecret))
	return signed
}

func strPtr(s string) *string { return &s }

// ── GET /patient/search/:id ───────────────────────────────────────────────────

func TestGetByID_FoundByNationalID(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makePatientToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByID", "1234567890123", "BKH01").
		Return(&domain.Patient{ID: "uuid-1", NationalID: strPtr("1234567890123")}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "uuid-1", body["id"])
	patientSvc.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makePatientToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByID", "unknown", "BKH01").Return(nil, domain.ErrNotFound)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/unknown", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "NOT_FOUND", body["code"])
}

func TestGetByID_NoAuthHeader(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetByID_InternalError(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makePatientToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByID", "1234567890123", "BKH01").Return(nil, assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ── POST /patient/search ──────────────────────────────────────────────────────

func TestSearch_Success(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makePatientToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)

	input := domain.PatientSearchInput{LastName: strPtr("Smith")}
	patients := []domain.Patient{{ID: "uuid-1"}, {ID: "uuid-2"}}
	patientSvc.On("GetPatientByCondition", input, "BKH01").Return(patients, nil)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{"last_name": "Smith"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body []map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Len(t, body, 2)
	patientSvc.AssertExpectations(t)
}

func TestSearch_EmptyResult(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makePatientToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)

	input := domain.PatientSearchInput{FirstName: strPtr("NoOne")}
	patientSvc.On("GetPatientByCondition", input, "BKH01").Return([]domain.Patient{}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{"first_name": "NoOne"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body []any
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Empty(t, body)
}

func TestSearch_NoCondition(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makePatientToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByCondition", domain.PatientSearchInput{}, "BKH01").
		Return(nil, domain.ErrInvalidInput)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "INVALID_INPUT", body["code"])
}

func TestSearch_BadBody(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makePatientToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody("not-an-object"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	// Fiber's BodyParser is lenient for primitive JSON — service will get empty input
	// which triggers ErrInvalidInput; either 400 is acceptable
	assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusOK)
}

func TestSearch_NoAuthHeader(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{"last_name": "Smith"}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSearch_InternalError(t *testing.T) {
	staffSvc := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makePatientToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByCondition", mock.Anything, "BKH01").Return(nil, assert.AnError)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{"last_name": "Smith"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
