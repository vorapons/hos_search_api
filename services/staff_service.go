package services

import (
	"errors"
	"regexp"
	"sync"
	"time"
	"unicode"

	"pt_search_hos/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// isStrongPassword requires: min 8 chars, at least 1 uppercase, 1 lowercase, 1 digit, 1 special char.
func isStrongPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

type staffClaims struct {
	Login      string `json:"login"`
	HospitalID string `json:"hospital_id"`
	Hospital   string `json:"hospital"`
	jwt.RegisteredClaims
}

type staffService struct {
	repo      domain.StaffRepository
	jwtSecret string
	blacklist map[string]struct{}
	mu        sync.RWMutex
}

func NewStaffService(repo domain.StaffRepository, jwtSecret string) domain.StaffService {
	return &staffService{
		repo:      repo,
		jwtSecret: jwtSecret,
		blacklist: make(map[string]struct{}),
	}
}

func (s *staffService) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", domain.ErrInvalidInput
	}

	staff, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if staff == nil {
		return "", domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(password)); err != nil {
		return "", domain.ErrUnauthorized
	}

	return s.generateToken(staff)
}

func (s *staffService) CreateStaff(email, password, hospitalName string) (string, error) {
	if email == "" || password == "" || hospitalName == "" {
		return "", domain.ErrInvalidInput
	}
	if !isValidEmail(email) {
		return "", domain.ErrInvalidInput
	}
	if !isStrongPassword(password) {
		return "", domain.ErrInvalidInput
	}

	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", domain.ErrStaffExists
	}

	hospital, err := s.repo.FindHospitalByName(hospitalName)
	if err != nil {
		return "", err
	}
	if hospital == nil {
		return "", domain.ErrHospitalNotFound
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	staff := &domain.Staff{
		HospitalID:   hospital.ID,
		HospitalName: hospital.Name,
		Email:        email,
		Password:     string(hashed),
	}

	if err := s.repo.Create(staff); err != nil {
		return "", err
	}

	return s.generateToken(staff)
}

func (s *staffService) Logout(tokenString string) error {
	expiresAt, err := s.tokenExpiry(tokenString)
	if err != nil {
		// If we can't parse the expiry, use 24h as safe default
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	if err := s.repo.AddBlacklistedToken(tokenString, expiresAt); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.blacklist[tokenString] = struct{}{}
	return nil
}

func (s *staffService) IsTokenBlacklisted(tokenString string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.blacklist[tokenString]
	return ok
}

// LoadBlacklist loads all non-expired blacklisted tokens from the DB into memory.
// Call this once on server startup.
func (s *staffService) LoadBlacklist() error {
	tokens, err := s.repo.LoadBlacklistedTokens()
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tokens {
		s.blacklist[t] = struct{}{}
	}
	return nil
}

func (s *staffService) generateToken(staff *domain.Staff) (string, error) {
	claims := staffClaims{
		Login:      staff.Email,
		HospitalID: staff.HospitalID,
		Hospital:   staff.HospitalName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *staffService) tokenExpiry(tokenString string) (time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &staffClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return time.Time{}, err
	}
	claims, ok := token.Claims.(*staffClaims)
	if !ok || claims.ExpiresAt == nil {
		return time.Time{}, errors.New("invalid claims")
	}
	return claims.ExpiresAt.Time, nil
}
