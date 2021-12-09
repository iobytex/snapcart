package gorm_errors

import (
	"errors"
	"gorm.io/gorm"
)

func GormError(gormError error)  error {
	switch  {
	case errors.Is(gormError,gorm.ErrInvalidData):
		return gorm.ErrInvalidData

	case errors.Is(gormError,gorm.ErrEmptySlice):
		return gorm.ErrEmptySlice

	case errors.Is(gormError,gorm.ErrInvalidField):
		return gorm.ErrInvalidField

	case errors.Is(gormError,gorm.ErrRecordNotFound):
		return gorm.ErrRecordNotFound
	default:
		return gormError
	}
}
