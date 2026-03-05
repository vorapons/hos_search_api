package domain

import "time"

// Hospital is a lightweight representation used across services.
type Hospital struct {
	ID   string
	Name string
}

// Staff represents an authenticated hospital staff member.
type Staff struct {
	ID           uint
	HospitalID   string
	HospitalName string
	Email        string
	Password     string  // bcrypt hash, never serialised
}

// StaffRepository is the port (interface) for the database adapter.
type StaffRepository interface {
	FindByEmail(email string) (*Staff, error)
	Create(staff *Staff) error
	FindHospitalByName(name string) (*Hospital, error)
	AddBlacklistedToken(token string, expiresAt time.Time) error
	LoadBlacklistedTokens() ([]string, error)
}

// StaffService is the port (interface) for the use-case layer.
type StaffService interface {
	Login(email, password, hospitalName string) (token string, err error)
	CreateStaff(email, password, hospitalName string) (token string, err error)
	Logout(token string) error
	IsTokenBlacklisted(token string) bool
	LoadBlacklist() error
}

// TokenClaims holds the data embedded in the JWT.
type TokenClaims struct {
	Login      string `json:"login"`
	HospitalID string `json:"hospital_id"`
	Hospital   string `json:"hospital"`
}
