package repository

import (
	"time"

	"gorm.io/gorm"
)

type HospitalModel struct {
	ID        string         `gorm:"type:varchar(5);primaryKey"`
	Name      string         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (HospitalModel) TableName() string { return "hospitals" }

type StaffModel struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	HospitalID string         `gorm:"type:varchar(5);not null"`
	Hospital   HospitalModel  `gorm:"foreignKey:HospitalID;references:ID"`
	Email      string         `gorm:"not null;uniqueIndex"`
	Password   string         `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (StaffModel) TableName() string { return "staff" }

type BlacklistedTokenModel struct {
	Token     string    `gorm:"primaryKey"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (BlacklistedTokenModel) TableName() string { return "blacklisted_tokens" }

type PatientModel struct {
	ID           string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	HospitalID   string         `gorm:"type:varchar(5);not null"`
	FirstNameTH  *string        `gorm:"column:first_name_th"`
	MiddleNameTH *string        `gorm:"column:middle_name_th"`
	LastNameTH   *string        `gorm:"column:last_name_th"`
	FirstNameEN  *string        `gorm:"column:first_name_en"`
	MiddleNameEN *string        `gorm:"column:middle_name_en"`
	LastNameEN   *string        `gorm:"column:last_name_en"`
	NationalID   *string        `gorm:"column:national_id"`
	PassportID   *string        `gorm:"column:passport_id"`
	PatientHN    *string        `gorm:"column:patient_hn"`
	DateOfBirth  *time.Time     `gorm:"column:date_of_birth"`
	Gender       *string        `gorm:"column:gender"`
	PhoneNumber  *string        `gorm:"column:phone_number"`
	Email        *string        `gorm:"column:email"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (PatientModel) TableName() string { return "patients" }
