package services

import (
	"strings"
	"sync"
	"time"

	"pt_search_hos/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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

	// DB requires at least one name; default to the part of the email before @
	nameEN := strings.SplitN(email, "@", 2)[0]
	staff := &domain.Staff{
		HospitalID:   hospital.ID,
		HospitalName: hospital.Name,
		Email:        email,
		Password:     string(hashed),
		NameEN:       &nameEN,
	}

	if err := s.repo.Create(staff); err != nil {
		return "", err
	}

	return s.generateToken(staff)
}

func (s *staffService) Logout(tokenString string) error {
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
