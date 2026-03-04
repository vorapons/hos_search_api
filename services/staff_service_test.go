package services_test

import (
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

// ── Login ─────────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(&domain.Staff{
		Email:    "user@example.com",
		Password: hashPassword("Pass1!xx"),
	}, nil)

	svc := newService(repo)
	token, err := svc.Login("user@example.com", "Pass1!xx")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

func TestLogin_EmptyEmail(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	_, err := svc.Login("", "Pass1!xx")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestLogin_EmptyPassword(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	_, err := svc.Login("user@example.com", "")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "ghost@example.com").Return(nil, nil)

	_, err := newService(repo).Login("ghost@example.com", "Pass1!xx")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(&domain.Staff{
		Email:    "user@example.com",
		Password: hashPassword("CorrectPass1!"),
	}, nil)

	_, err := newService(repo).Login("user@example.com", "WrongPass1!")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestLogin_DBError(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(nil, assert.AnError)

	_, err := newService(repo).Login("user@example.com", "Pass1!xx")
	assert.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrUnauthorized)
}

// ── CreateStaff ───────────────────────────────────────────────────────────────

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

func TestCreateStaff_EmptyFields(t *testing.T) {
	svc := newService(new(mockStaffRepo))

	_, err := svc.CreateStaff("", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)

	_, err = svc.CreateStaff("user@example.com", "", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)

	_, err = svc.CreateStaff("user@example.com", "Pass1!xx", "")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestCreateStaff_InvalidEmail(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	_, err := svc.CreateStaff("not-an-email", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
}

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

func TestCreateStaff_AlreadyExists(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "dup@example.com").Return(&domain.Staff{
		Email: "dup@example.com",
	}, nil)

	_, err := newService(repo).CreateStaff("dup@example.com", "Pass1!xx", "Bangkok Hospital")
	assert.ErrorIs(t, err, domain.ErrStaffExists)
}

func TestCreateStaff_HospitalNotFound(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("FindByEmail", "user@example.com").Return(nil, nil)
	repo.On("FindHospitalByName", "Unknown Hospital").Return(nil, nil)

	_, err := newService(repo).CreateStaff("user@example.com", "Pass1!xx", "Unknown Hospital")
	assert.ErrorIs(t, err, domain.ErrHospitalNotFound)
}

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

func TestLogout_InvalidTokenFallsBackTo24h(t *testing.T) {
	repo := new(mockStaffRepo)
	// Can't parse expiry from "badtoken" — service should fall back to 24h default
	repo.On("AddBlacklistedToken", "badtoken", mock.AnythingOfType("time.Time")).Return(nil)

	svc := newService(repo)
	err := svc.Logout("badtoken")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLogout_DBError(t *testing.T) {
	repo := new(mockStaffRepo)
	token := makeSignedToken("user@example.com", "BKH01", "Bangkok Hospital")
	repo.On("AddBlacklistedToken", token, mock.AnythingOfType("time.Time")).Return(assert.AnError)

	err := newService(repo).Logout(token)
	assert.Error(t, err)
}

// ── IsTokenBlacklisted ────────────────────────────────────────────────────────

func TestIsTokenBlacklisted_FalseInitially(t *testing.T) {
	svc := newService(new(mockStaffRepo))
	assert.False(t, svc.IsTokenBlacklisted("some-token"))
}

// ── LoadBlacklist ─────────────────────────────────────────────────────────────

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

func TestLoadBlacklist_DBError(t *testing.T) {
	repo := new(mockStaffRepo)
	repo.On("LoadBlacklistedTokens").Return(nil, assert.AnError)

	err := newService(repo).LoadBlacklist()
	assert.Error(t, err)
}
