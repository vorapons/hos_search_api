package domain

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
}

// StaffService is the port (interface) for the use-case layer.
type StaffService interface {
	Login(email, password string) (token string, err error)
	CreateStaff(email, password, hospitalName string) (token string, err error)
	Logout(token string) error
	IsTokenBlacklisted(token string) bool
}

// TokenClaims holds the data embedded in the JWT.
type TokenClaims struct {
	Login      string `json:"login"`
	HospitalID string `json:"hospital_id"`
	Hospital   string `json:"hospital"`
}
