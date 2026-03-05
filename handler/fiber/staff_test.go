package fiberhandler_test

import (
	"bytes"
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

const testJWTSecret = "test-secret"

// ── mock ──────────────────────────────────────────────────────────────────────

type mockStaffService struct {
	mock.Mock
}

func (m *mockStaffService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *mockStaffService) CreateStaff(email, password, hospital string) (string, error) {
	args := m.Called(email, password, hospital)
	return args.String(0), args.Error(1)
}

func (m *mockStaffService) Logout(token string) error {
	return m.Called(token).Error(0)
}

func (m *mockStaffService) IsTokenBlacklisted(token string) bool {
	return m.Called(token).Bool(0)
}

func (m *mockStaffService) LoadBlacklist() error {
	return m.Called().Error(0)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func setupApp(svc domain.StaffService) *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	staffH   := fiberhandler.NewStaffHandler(svc)
	patientH := fiberhandler.NewPatientHandler(nil)
	fiberhandler.SetupRoutes(app, staffH, patientH, testJWTSecret, svc.IsTokenBlacklisted)
	return app
}

// makeToken generates a valid signed JWT for use in tests.
func makeToken(login, hospitalID, hospital string) string {
	type claims struct {
		Login      string `json:"login"`
		HospitalID string `json:"hospital_id"`
		Hospital   string `json:"hospital"`
		jwt.RegisteredClaims
	}
	c := claims{
		Login:      login,
		HospitalID: hospitalID,
		Hospital:   hospital,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, _ := token.SignedString([]byte(testJWTSecret))
	return signed
}

func jsonBody(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

// ── GET /hello ────────────────────────────────────────────────────────────────

// positive: health check → 200 with status and timestamp
func TestHello_HealthCheck(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "ok", body["status"])
	assert.NotEmpty(t, body["timestamp"])
}

// ── POST /staff/login ─────────────────────────────────────────────────────────

// positive: valid credentials → 200 with token
func TestLogin_Success(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "user@example.com", "Pass1!xx").Return("tok123", nil)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "user@example.com", "password": "Pass1!xx"}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "tok123", body["token"])
	svc.AssertExpectations(t)
}

// negative: malformed JSON body → 400
func TestLogin_BadBody(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// negative: empty email/password → 400 INVALID_INPUT
func TestLogin_InvalidInput(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "", "").Return("", domain.ErrInvalidInput)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "", "password": ""}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "INVALID_INPUT", body["code"])
}

// negative: wrong password → 401 UNAUTHORIZED
func TestLogin_Unauthorized(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "user@example.com", "WrongPass1!").Return("", domain.ErrUnauthorized)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "user@example.com", "password": "WrongPass1!"}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "UNAUTHORIZED", body["code"])
}

// negative: login value is not an email address → 401 UNAUTHORIZED (service treats it as user not found)
func TestLogin_InvalidEmailFormat(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "notanemail", "Pass1!xx").Return("", domain.ErrUnauthorized)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "notanemail", "password": "Pass1!xx"}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "UNAUTHORIZED", body["code"])
}

// negative: unexpected service error → 500
func TestLogin_InternalError(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "user@example.com", "Pass1!xx").Return("", assert.AnError)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "user@example.com", "password": "Pass1!xx"}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ── POST /staff/create ────────────────────────────────────────────────────────

// positive: new staff account created → 201 with token
func TestCreate_Success(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "new@example.com", "Pass1!xx", "Bangkok Hospital").Return("tok456", nil)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{
			"login":    "new@example.com",
			"password": "Pass1!xx",
			"hospital": "Bangkok Hospital",
		}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "tok456", body["token"])
	svc.AssertExpectations(t)
}

// negative: malformed JSON body → 400
func TestCreate_BadBody(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// negative: empty fields → 400 INVALID_INPUT
func TestCreate_InvalidInput(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "", "", "").Return("", domain.ErrInvalidInput)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{"login": "", "password": "", "hospital": ""}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "INVALID_INPUT", body["code"])
}

// negative: login is not a valid email address → 400 INVALID_INPUT
func TestCreate_InvalidEmailFormat(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "notanemail", "Pass1!xx", "Bangkok Hospital").Return("", domain.ErrInvalidInput)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{
			"login":    "notanemail",
			"password": "Pass1!xx",
			"hospital": "Bangkok Hospital",
		}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "INVALID_INPUT", body["code"])
}

// negative: password does not meet strength requirements → 400 INVALID_INPUT
func TestCreate_WeakPassword(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "new@example.com", "weak", "Bangkok Hospital").Return("", domain.ErrInvalidInput)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{
			"login":    "new@example.com",
			"password": "weak",
			"hospital": "Bangkok Hospital",
		}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "INVALID_INPUT", body["code"])
}

// negative: email already registered → 409 CONFLICT
func TestCreate_StaffExists(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "dup@example.com", "Pass1!xx", "Bangkok Hospital").Return("", domain.ErrStaffExists)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{
			"login":    "dup@example.com",
			"password": "Pass1!xx",
			"hospital": "Bangkok Hospital",
		}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "CONFLICT", body["code"])
}

// negative: hospital not in DB → 404 NOT_FOUND
func TestCreate_HospitalNotFound(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "user@example.com", "Pass1!xx", "Unknown Hospital").Return("", domain.ErrHospitalNotFound)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{
			"login":    "user@example.com",
			"password": "Pass1!xx",
			"hospital": "Unknown Hospital",
		}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "NOT_FOUND", body["code"])
}

// negative: hospital name with wrong casing → 404 NOT_FOUND (lookup is case-sensitive)
func TestCreate_HospitalNameWrongCase(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "user@example.com", "Pass1!xx", "bangkok hospital").Return("", domain.ErrHospitalNotFound)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{
			"login":    "user@example.com",
			"password": "Pass1!xx",
			"hospital": "bangkok hospital",
		}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "NOT_FOUND", body["code"])
}

// negative: unexpected service error → 500
func TestCreate_InternalError(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "user@example.com", "Pass1!xx", "Bangkok Hospital").Return("", assert.AnError)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{
			"login":    "user@example.com",
			"password": "Pass1!xx",
			"hospital": "Bangkok Hospital",
		}))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ── GET /staff/hello ──────────────────────────────────────────────────────────

// positive: valid token → 200 with login info
func TestHello_Success(t *testing.T) {
	svc := new(mockStaffService)
	token := makeToken("user@example.com", "BKH01", "Bangkok Hospital")
	svc.On("IsTokenBlacklisted", token).Return(false)

	req, _ := http.NewRequest(http.MethodGet, "/staff/hello", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "user@example.com", body["login"])
	assert.Equal(t, "Bangkok Hospital", body["hospital"])
	assert.NotEmpty(t, body["expires_at"])
}

// negative: missing Authorization header → 401 Unauthorized
func TestHello_NoAuthHeader(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodGet, "/staff/hello", nil)

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// negative: malformed/invalid JWT → 401 Unauthorized
func TestHello_InvalidToken(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("IsTokenBlacklisted", "bad.token.here").Return(false)

	req, _ := http.NewRequest(http.MethodGet, "/staff/hello", nil)
	req.Header.Set("Authorization", "Bearer bad.token.here")

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// negative: revoked/blacklisted token → 401 UNAUTHORIZED
func TestHello_BlacklistedToken(t *testing.T) {
	svc := new(mockStaffService)
	token := makeToken("user@example.com", "BKH01", "Bangkok Hospital")
	svc.On("IsTokenBlacklisted", token).Return(true)

	req, _ := http.NewRequest(http.MethodGet, "/staff/hello", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "UNAUTHORIZED", body["code"])
}

// ── GET /staff/logout ─────────────────────────────────────────────────────────

// positive: valid token, logout succeeds → 200
func TestLogout_Success(t *testing.T) {
	svc := new(mockStaffService)
	token := makeToken("user@example.com", "BKH01", "Bangkok Hospital")
	svc.On("IsTokenBlacklisted", token).Return(false)
	svc.On("Logout", token).Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/staff/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "Logged out successfully", body["message"])
	svc.AssertExpectations(t)
}

// negative: missing Authorization header → 401 Unauthorized
func TestLogout_NoAuthHeader(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodGet, "/staff/logout", nil)

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// negative: service fails to blacklist token → 500
func TestLogout_ServiceError(t *testing.T) {
	svc := new(mockStaffService)
	token := makeToken("user@example.com", "BKH01", "Bangkok Hospital")
	svc.On("IsTokenBlacklisted", token).Return(false)
	svc.On("Logout", token).Return(assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/staff/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := setupApp(svc).Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
