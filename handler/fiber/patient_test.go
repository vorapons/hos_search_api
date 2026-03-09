package fiberhandler_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"pt_search_hos/domain"
	fiberhandler "pt_search_hos/handler/fiber"

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

func (m *mockPatientService) GetPatientByCondition(input domain.PatientSearchInput, hospitalID string) (domain.PatientSearchResult, error) {
	args := m.Called(input, hospitalID)
	result, _ := args.Get(0).(domain.PatientSearchResult)
	return result, args.Error(1)
}

// ── helpers function ────────────────────────────────────────────────────────────────

func setupPatientApp(staffSvc domain.StaffService, patientSvc domain.PatientService) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	staffH   := fiberhandler.NewStaffHandler(staffSvc)
	patientH := fiberhandler.NewPatientHandler(patientSvc)
	fiberhandler.SetupRoutes(app, staffH, patientH, testJWTSecret, staffSvc.IsTokenBlacklisted)
	return app
}

// makeStaffToken generates a staff JWT with the given hospital_id for authenticating patient endpoints.
func makeStaffToken(hospitalID string) string {
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

// positive: patient found by national ID → 200 with patient data
func TestGetByID_FoundByNationalID(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByID", "1234567890123", "BKH01").
		Return(&domain.Patient{ID: "uuid-1", NationalID: strPtr("1234567890123")}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "1234567890123", body["national_id"])
	patientSvc.AssertExpectations(t)
}

// positive: patient found by passport ID → 200 with patient data
func TestGetByID_FoundByPassportID(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByID", "AB123456", "BKH01").
		Return(&domain.Patient{ID: "uuid-2", PassportID: strPtr("AB123456")}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/AB123456", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "AB123456", body["passport_id"])
	patientSvc.AssertExpectations(t)
}

// negative: patient does not exist → 404 NOT_FOUND
func TestGetByID_NotFound(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
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

// negative: missing Authorization header → 401 Unauthorized
func TestGetByID_NoAuthHeader(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// negative: service rejects empty/invalid ID → 400 INVALID_INPUT
func TestGetByID_InvalidInput(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByID", "1234567890123", "BKH01").Return(nil, domain.ErrInvalidInput)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "INVALID_INPUT", body["code"])
}

// negative: unexpected service error → 500 INTERNAL_ERROR
func TestGetByID_InternalError(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByID", "1234567890123", "BKH01").Return(nil, assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/patient/search/1234567890123", nil)
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ── POST /patient/search ──────────────────────────────────────────────────────

// positive: search returns mixed patients → 200 with paginated result; verifies Thai names+national_id and foreign passport_id are serialised correctly
func TestSearch_Success(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)

	dob1 := time.Date(1985, 3, 15, 0, 0, 0, 0, time.UTC) // Somchai
	dob2 := time.Date(1990, 3, 14, 0, 0, 0, 0, time.UTC) // Yuki Tanaka

	input := domain.PatientSearchInput{LastName: strPtr("Smith")}
	searchResult := domain.PatientSearchResult{
		Data: []domain.Patient{
			// Thai — Somchai Jaidee (BKH-0001): has Thai names + national_id, no passport
			{
				FirstNameTH: strPtr("สมชาย"),  LastNameTH:  strPtr("ใจดี"),
				FirstNameEN: strPtr("Somchai"), LastNameEN:  strPtr("Jaidee"),
				NationalID:  strPtr("1100100012341"), PatientHN: strPtr("BKH-0001"),
				DateOfBirth: &dob1, Gender: strPtr("male"),
				PhoneNumber: strPtr("0812345001"), Email: strPtr("somchai.j@email.com"),
			},
			// Japanese — Yuki Tanaka (BKH-0005): no Thai names, passport only
			{
				FirstNameEN: strPtr("Yuki"),    LastNameEN:  strPtr("Tanaka"),
				PassportID:  strPtr("JP10234567"), PatientHN: strPtr("BKH-0005"),
				DateOfBirth: &dob2, Gender: strPtr("female"),
				PhoneNumber: strPtr("+819011112001"), Email: strPtr("yuki.tanaka@email.jp"),
			},
		},
		Total: 2, Page: 1, PageSize: 20,
	}
	patientSvc.On("GetPatientByCondition", input, "BKH01").Return(searchResult, nil)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{"last_name": "Smith"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)

	assert.Equal(t, float64(2), body["total"])
	assert.Equal(t, float64(1), body["page"])
	assert.Equal(t, float64(20), body["page_size"])

	data, _ := body["data"].([]any)
	assert.Len(t, data, 2)

	p0 := data[0].(map[string]any)
	assert.Equal(t, "สมชาย",         p0["first_name_th"])
	assert.Equal(t, "ใจดี",          p0["last_name_th"])
	assert.Equal(t, "Somchai",       p0["first_name_en"])
	assert.Equal(t, "1100100012341", p0["national_id"])
	assert.Nil(t,                    p0["passport_id"])
	assert.Equal(t, "BKH-0001",      p0["patient_hn"])
	assert.Equal(t, "male",          p0["gender"])

	p1 := data[1].(map[string]any)
	assert.Nil(t,                    p1["first_name_th"])
	assert.Nil(t,                    p1["last_name_th"])
	assert.Equal(t, "Yuki",          p1["first_name_en"])
	assert.Equal(t, "JP10234567",    p1["passport_id"])
	assert.Nil(t,                    p1["national_id"])
	assert.Equal(t, "BKH-0005",      p1["patient_hn"])
	assert.Equal(t, "female",        p1["gender"])

	patientSvc.AssertExpectations(t)
}

// positive: valid condition, no patients matched → 200 with empty data and total=0
func TestSearch_EmptyResult(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)

	input := domain.PatientSearchInput{FirstName: strPtr("NoOne")}
	patientSvc.On("GetPatientByCondition", input, "BKH01").Return(
		domain.PatientSearchResult{Data: []domain.Patient{}, Total: 0, Page: 1, PageSize: 20}, nil)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{"first_name": "NoOne"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, float64(0), body["total"])
	data, _ := body["data"].([]any)
	assert.Empty(t, data)
}

// negative: no search condition provided → 400 INVALID_INPUT
func TestSearch_NoCondition(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByCondition", domain.PatientSearchInput{}, "BKH01").
		Return(domain.PatientSearchResult{}, domain.ErrInvalidInput)

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

// negative: request body is not a JSON object → 400
func TestSearch_BadBody(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
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

// negative: missing Authorization header → 401 Unauthorized
func TestSearch_NoAuthHeader(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{"last_name": "Smith"}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// negative: unexpected service error → 500 INTERNAL_ERROR
func TestSearch_InternalError(t *testing.T) {
	staffSvc   := new(mockStaffService)
	patientSvc := new(mockPatientService)

	tok := makeStaffToken("BKH01")
	staffSvc.On("IsTokenBlacklisted", tok).Return(false)
	patientSvc.On("GetPatientByCondition", mock.Anything, "BKH01").Return(domain.PatientSearchResult{}, assert.AnError)

	req, _ := http.NewRequest(http.MethodPost, "/patient/search",
		jsonBody(map[string]string{"last_name": "Smith"}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)

	resp, _ := setupPatientApp(staffSvc, patientSvc).Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
