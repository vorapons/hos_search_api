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
