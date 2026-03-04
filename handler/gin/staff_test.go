package ginhandler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pt_search_hos/domain"
	ginhandler "pt_search_hos/handler/gin"

	"github.com/gin-gonic/gin"
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

func setupRouter(svc domain.StaffService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	staffH   := ginhandler.NewStaffHandler(svc)
	patientH := ginhandler.NewPatientHandler(nil)
	ginhandler.SetupRoutes(r, staffH, patientH, testJWTSecret, svc.IsTokenBlacklisted)
	return r
}

func perform(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
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

// ── POST /staff/login ─────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "user@example.com", "Pass1!xx").Return("tok123", nil)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "user@example.com", "password": "Pass1!xx"}))
	req.Header.Set("Content-Type", "application/json")

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "tok123", body["token"])
	svc.AssertExpectations(t)
}

func TestLogin_BadBody(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_InvalidInput(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "", "").Return("", domain.ErrInvalidInput)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "", "password": ""}))
	req.Header.Set("Content-Type", "application/json")

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "INVALID_INPUT", body["code"])
}

func TestLogin_Unauthorized(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "user@example.com", "WrongPass1!").Return("", domain.ErrUnauthorized)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "user@example.com", "password": "WrongPass1!"}))
	req.Header.Set("Content-Type", "application/json")

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "UNAUTHORIZED", body["code"])
}

func TestLogin_InternalError(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("Login", "user@example.com", "Pass1!xx").Return("", assert.AnError)

	req, _ := http.NewRequest(http.MethodPost, "/staff/login",
		jsonBody(map[string]string{"login": "user@example.com", "password": "Pass1!xx"}))
	req.Header.Set("Content-Type", "application/json")

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── POST /staff/create ────────────────────────────────────────────────────────

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

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "tok456", body["token"])
	svc.AssertExpectations(t)
}

func TestCreate_BadBody(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_InvalidInput(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("CreateStaff", "", "", "").Return("", domain.ErrInvalidInput)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create",
		jsonBody(map[string]string{"login": "", "password": "", "hospital": ""}))
	req.Header.Set("Content-Type", "application/json")

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "INVALID_INPUT", body["code"])
}

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

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusConflict, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "CONFLICT", body["code"])
}

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

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "NOT_FOUND", body["code"])
}

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

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── GET /staff/hello ──────────────────────────────────────────────────────────

func TestHello_Success(t *testing.T) {
	svc := new(mockStaffService)
	token := makeToken("user@example.com", "BKH01", "Bangkok Hospital")
	svc.On("IsTokenBlacklisted", token).Return(false)

	req, _ := http.NewRequest(http.MethodGet, "/staff/hello", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "user@example.com", body["login"])
	assert.Equal(t, "Bangkok Hospital", body["hospital"])
	assert.NotEmpty(t, body["expires_at"])
}

func TestHello_NoAuthHeader(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodGet, "/staff/hello", nil)

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHello_InvalidToken(t *testing.T) {
	svc := new(mockStaffService)
	svc.On("IsTokenBlacklisted", "bad.token.here").Return(false)

	req, _ := http.NewRequest(http.MethodGet, "/staff/hello", nil)
	req.Header.Set("Authorization", "Bearer bad.token.here")

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHello_BlacklistedToken(t *testing.T) {
	svc := new(mockStaffService)
	token := makeToken("user@example.com", "BKH01", "Bangkok Hospital")
	svc.On("IsTokenBlacklisted", token).Return(true)

	req, _ := http.NewRequest(http.MethodGet, "/staff/hello", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "UNAUTHORIZED", body["code"])
}

// ── GET /staff/logout ─────────────────────────────────────────────────────────

func TestLogout_Success(t *testing.T) {
	svc := new(mockStaffService)
	token := makeToken("user@example.com", "BKH01", "Bangkok Hospital")
	svc.On("IsTokenBlacklisted", token).Return(false)
	svc.On("Logout", token).Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/staff/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	assert.Equal(t, "Logged out successfully", body["message"])
	svc.AssertExpectations(t)
}

func TestLogout_NoAuthHeader(t *testing.T) {
	svc := new(mockStaffService)
	req, _ := http.NewRequest(http.MethodGet, "/staff/logout", nil)

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogout_ServiceError(t *testing.T) {
	svc := new(mockStaffService)
	token := makeToken("user@example.com", "BKH01", "Bangkok Hospital")
	svc.On("IsTokenBlacklisted", token).Return(false)
	svc.On("Logout", token).Return(assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/staff/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := perform(setupRouter(svc), req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
