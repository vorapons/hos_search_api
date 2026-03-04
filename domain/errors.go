package domain

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("conflict")
	ErrHospitalNotFound = errors.New("hospital not found")
	ErrStaffExists      = errors.New("staff already exists")
)
