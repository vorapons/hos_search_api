package services_test

import (
	"strings"
	"testing"
	"time"

	"pt_search_hos/domain"
	"pt_search_hos/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const testSecret = "test-jwt-secret"

// ── mock repo ─────────────────────────────────────────────────────────────────

type mockStaffRepo struct {
	mock.Mock
}

func (m *mockStaffRepo) FindByEmail(email string) (*domain.Staff, error) {
	args := m.Called(email)
	staff, _ := args.Get(0).(*domain.Staff)
	return staff, args.Error(1)
}

func (m *mockStaffRepo) Create(staff *domain.Staff) error {
	return m.Called(staff).Error(0)
}

func (m *mockStaffRepo) FindHospitalByName(name string) (*domain.Hospital, error) {
	args := m.Called(name)
	hospital, _ := args.Get(0).(*domain.Hospital)
	return hospital, args.Error(1)
}

func (m *mockStaffRepo) AddBlacklistedToken(token string, expiresAt time.Time) error {
	return m.Called(token, expiresAt).Error(0)
}

func (m *mockStaffRepo) LoadBlacklistedTokens() ([]string, error) {
	args := m.Called()
	tokens, _ := args.Get(0).([]string)
	return tokens, args.Error(1)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func newService(repo domain.StaffRepository) domain.StaffService {
	return services.NewStaffService(repo, testSecret)
}

func hashPassword(pw string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	return string(h)
}

// makeSignedToken creates a real JWT with the test secret for use in Logout tests.
func makeSignedToken(email, hospitalID, hospital string) string {
	type claims struct {
		Login      string `json:"login"`
		HospitalID string `json:"hospital_id"`
		Hospital   string `json:"hospital"`
		jwt.RegisteredClaims
	}
	c := claims{
		Login:      email,
		HospitalID: hospitalID,
		Hospital:   hospital,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, _ := token.SignedString([]byte(testSecret))
	return signed
}

// ! ── Login ─────────────────────────────────────────────────────────────────────

// positive: valid credentials → JWT token returned
func TestLogin_Success(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(&domain.Staff{
		Email:        "user@example.com",
		Password:     hashPassword("Pass1!xx"),
		HospitalName: "Bangkok Hospital",
	}, nil)

	svc := newService(repo)
	token, err := svc.Login("user@example.com", "Pass1!xx", "Bangkok Hospital")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

// negative: empty email → ErrInvalidInput
func TestLogin_EmptyEmail(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	_, err := svc.Login("", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// negative: empty password → ErrInvalidInput
func TestLogin_EmptyPassword(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	_, err := svc.Login("user@example.com", "", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// negative: empty hospital → ErrInvalidInput
func TestLogin_EmptyHospital(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	_, err := svc.Login("user@example.com", "Pass1!xx", "")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// negative: login is not an email format → ErrUnauthorized (user enumeration prevention — same as not found)
func TestLogin_NotEmailFormat(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "notanemail").Return(nil, nil)

	_, err := newService(repo).Login("notanemail", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
	repo.AssertExpectations(t) // repo IS called — no short-circuit on format
}

// negative: email not in DB → ErrUnauthorized
func TestLogin_UserNotFound(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "ghost@example.com").Return(nil, nil)

	_, err := newService(repo).Login("ghost@example.com", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

// negative: wrong password → ErrUnauthorized
func TestLogin_WrongPassword(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(&domain.Staff{
		Email:        "user@example.com",
		Password:     hashPassword("CorrectPass1!"),
		HospitalName: "Bangkok Hospital",
	}, nil)

	_, err := newService(repo).Login("user@example.com", "WrongPass1!", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

// negative: hospital name does not match staff's hospital → ErrUnauthorized
func TestLogin_WrongHospital(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(&domain.Staff{
		Email:        "user@example.com",
		Password:     hashPassword("Pass1!xx"),
		HospitalName: "Bangkok Hospital",
	}, nil)

	_, err := newService(repo).Login("user@example.com", "Pass1!xx", "Other Hospital")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

// negative: repository returns DB error → propagated as-is
func TestLogin_DBError(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(nil, assert.AnError)

	_, err := newService(repo).Login("user@example.com", "Pass1!xx", "Bangkok Hospital")
	assert.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrUnauthorized)
}

// ! ── CreateStaff ───────────────────────────────────────────────────────────────

// positive: new staff created, token returned
func TestCreateStaff_Success(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "new@example.com").Return(nil, nil)
	repo.On("FindHospitalByName", "Bangkok Hospital").Return(&domain.Hospital{
		ID: "BKH01", Name: "Bangkok Hospital",
	}, nil)
	repo.On("Create", mock.AnythingOfType("*domain.Staff")).Return(nil)

	svc := newService(repo)
	token, err := svc.CreateStaff("new@example.com", "Pass1!xx", "Bangkok Hospital")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

// negative: any field empty → ErrInvalidInput
func TestCreateStaff_EmptyFields(t *testing.T) {
	svc := newService(new(mockStaffRepo))

	_, err := svc.CreateStaff("", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)

	_, err = svc.CreateStaff("user@example.com", "", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)

	_, err = svc.CreateStaff("user@example.com", "Pass1!xx", "")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// negative: malformed email → ErrInvalidInput
func TestCreateStaff_InvalidEmail(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	_, err := svc.CreateStaff("not-an-email", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

// negative: password fails strength check → ErrInvalidInput
func TestCreateStaff_WeakPassword(t *testing.T) {
	svc := newService(new(mockStaffRepo))

	cases := []string{
		"short",        // too short
		"alllowercase1!", // no uppercase
		"ALLUPPERCASE1!", // no lowercase
		"NoSpecialChar1", // no special char
		"NoDigit!xx",     // no digit
	}
	for _, pw := range cases {
		_, err := svc.CreateStaff("user@example.com", pw, "Bangkok Hospital")
		assert.ErrorIs(t, err, domain.ErrInvalidInput, "expected weak password to fail: %s", pw)
	}
}

// negative: email already registered → ErrStaffExists
func TestCreateStaff_AlreadyExists(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "dup@example.com").Return(&domain.Staff{
		Email: "dup@example.com",
	}, nil)

	_, err := newService(repo).CreateStaff("dup@example.com", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrStaffExists)
}

// negative: hospital not in DB → ErrHospitalNotFound
func TestCreateStaff_HospitalNotFound(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(nil, nil)
	repo.On("FindHospitalByName", "Unknown Hospital").Return(nil, nil)

	_, err := newService(repo).CreateStaff("user@example.com", "Pass1!xx", "Unknown Hospital")
	assert.ErrorIs(t, err, domain.ErrHospitalNotFound)
}

// negative: DB error during duplicate email check → propagated as-is
func TestCreateStaff_FindByEmailDBError(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(nil, assert.AnError)

	_, err := newService(repo).CreateStaff("user@example.com", "Pass1!xx", "Bangkok Hospital")
	assert.Error(t, err)
}

// negative: DB error during hospital lookup → propagated as-is
func TestCreateStaff_FindHospitalDBError(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(nil, nil)
	repo.On("FindHospitalByName", "Bangkok Hospital").Return(nil, assert.AnError)

	_, err := newService(repo).CreateStaff("user@example.com", "Pass1!xx", "Bangkok Hospital")
	assert.Error(t, err)
}

// negative: password > 72 bytes passes strength check but bcrypt rejects it → error
func TestCreateStaff_PasswordTooLong(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(nil, nil)
	repo.On("FindHospitalByName", "Bangkok Hospital").Return(&domain.Hospital{
		ID: "BKH01", Name: "Bangkok Hospital",
	}, nil)

	// 73-byte password: passes isStrongPassword but bcrypt rejects it
	longPassword := "Pass1!" + strings.Repeat("a", 67)
	_, err := newService(repo).CreateStaff("user@example.com", longPassword, "Bangkok Hospital")
	assert.Error(t, err)
}

// negative: DB error during staff creation → propagated as-is
func TestCreateStaff_DBCreateError(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(nil, nil)
	repo.On("FindHospitalByName", "Bangkok Hospital").Return(&domain.Hospital{
		ID: "BKH01", Name: "Bangkok Hospital",
	}, nil)
	repo.On("Create", mock.AnythingOfType("*domain.Staff")).Return(assert.AnError)

	_, err := newService(repo).CreateStaff("user@example.com", "Pass1!xx", "Bangkok Hospital")
	assert.Error(t, err)
}

// ── Logout ────────────────────────────────────────────────────────────────────

// positive: token blacklisted in DB and in-memory cache
func TestLogout_PersistsToDBAndMemory(t *testing.T) {
	repo := new(mockStaffRepo)
	token := makeSignedToken("user@example.com", "BKH01", "Bangkok Hospital")
	repo.On("AddBlacklistedToken", token, mock.AnythingOfType("time.Time")).Return(nil)

	svc := newService(repo)
	err := svc.Logout(token)

	assert.NoError(t, err)
	assert.True(t, svc.IsTokenBlacklisted(token))
	repo.AssertExpectations(t)
}

// negative (edge): unparseable token string → falls back to 24h expiry, still succeeds
func TestLogout_InvalidTokenFallsBackTo24h(t *testing.T) {
	repo := new(mockStaffRepo)
	// Can't parse expiry from "badtoken" — service should fall back to 24h default
	repo.On("AddBlacklistedToken", "badtoken", mock.AnythingOfType("time.Time")).Return(nil)

	svc := newService(repo)
	err := svc.Logout("badtoken")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// negative (edge): valid JWT without exp claim → falls back to 24h expiry, still succeeds
func TestLogout_TokenWithNoExpiryClaim(t *testing.T) {
	repo := new(mockStaffRepo)

	// Valid JWT signed with correct secret but no exp claim → claims.ExpiresAt == nil (line 193)
	type noExpClaims struct {
		Login string `json:"login"`
		jwt.RegisteredClaims
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, noExpClaims{Login: "user@example.com"})
	tokenString, _ := tok.SignedString([]byte(testSecret))

	// tokenExpiry returns "invalid claims" → Logout falls back to 24h default
	repo.On("AddBlacklistedToken", tokenString, mock.AnythingOfType("time.Time")).Return(nil)

	err := newService(repo).Logout(tokenString)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// negative: DB error when blacklisting token → error returned
func TestLogout_DBError(t *testing.T) {
	repo := new(mockStaffRepo)
	token := makeSignedToken("user@example.com", "BKH01", "Bangkok Hospital")
	repo.On("AddBlacklistedToken", token, mock.AnythingOfType("time.Time")).Return(assert.AnError)

	err := newService(repo).Logout(token)
	assert.Error(t, err)
}

// ── IsTokenBlacklisted ────────────────────────────────────────────────────────

// positive: freshly created service has empty blacklist → always false
func TestIsTokenBlacklisted_FalseInitially(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	assert.False(t, svc.IsTokenBlacklisted("some-token"))
}

// ── LoadBlacklist ─────────────────────────────────────────────────────────────

// positive: DB tokens loaded into in-memory blacklist on startup
func TestLoadBlacklist_LoadsTokensIntoMemory(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("LoadBlacklistedTokens").Return([]string{"tok1", "tok2"}, nil)

	svc := newService(repo)
	err := svc.LoadBlacklist()

	assert.NoError(t, err)
	assert.True(t, svc.IsTokenBlacklisted("tok1"))
	assert.True(t, svc.IsTokenBlacklisted("tok2"))
	assert.False(t, svc.IsTokenBlacklisted("tok3"))
	repo.AssertExpectations(t)
}

// negative: DB error during load → error returned
func TestLoadBlacklist_DBError(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("LoadBlacklistedTokens").Return(nil, assert.AnError)

	err := newService(repo).LoadBlacklist()
	assert.Error(t, err)
}
