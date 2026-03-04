package repository

import (
	"time"

	"gorm.io/gorm"
)

type HospitalModel struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	Name      string         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (HospitalModel) TableName() string { return "hospitals" }

type StaffModel struct {
	ID         string         `gorm:"type:uuid;primaryKey"`
	HospitalID string         `gorm:"type:uuid;not null"`
	Hospital   HospitalModel  `gorm:"foreignKey:HospitalID;references:ID"`
	Email      string         `gorm:"not null;uniqueIndex"`
	Password   string         `gorm:"not null"`
	Role       string         `gorm:"not null;default:'staff'"`
	NameTH     *string        `gorm:"column:name_th"`
	NameEN     *string        `gorm:"column:name_en"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (StaffModel) TableName() string { return "staff" }
