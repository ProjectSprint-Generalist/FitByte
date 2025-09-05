package models

import (
	"errors"
	"time"
)

// Enums for validation
const (
	PreferenceCardio = "CARDIO"
	PreferenceWeight = "WEIGHT"

	WeightUnitKG  = "KG"
	WeightUnitLBS = "LBS"

	HeightUnitCM   = "CM"
	HeightUnitINCH = "INCH"
)

// User represents a user in the system
type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"uniqueIndex;not null"`
	Password   string    `json:"-" gorm:"not null"`
	Name       *string   `json:"name" gorm:"type:varchar(255)"`
	Preference *string   `json:"preference" gorm:"type:varchar(20);check:preference IN ('CARDIO', 'WEIGHT')"`
	WeightUnit *string   `json:"weightUnit" gorm:"type:varchar(10);check:weight_unit IN ('KG', 'LBS')"`
	HeightUnit *string   `json:"heightUnit" gorm:"type:varchar(10);check:height_unit IN ('CM', 'INCH')"`
	Weight     *float64  `json:"weight" gorm:"type:decimal(6,2)"`
	Height     *float64  `json:"height" gorm:"type:decimal(6,2)"`
	ImageURI   *string   `json:"imageUri" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email      string   `json:"email" binding:"required,email"`
	Name       *string  `json:"name"`
	Preference *string  `json:"preference"`
	WeightUnit *string  `json:"weightUnit"`
	HeightUnit *string  `json:"heightUnit"`
	Weight     *float64 `json:"weight"`
	Height     *float64 `json:"height"`
	ImageURI   *string  `json:"imageUri"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email      *string  `json:"email,omitempty" binding:"omitempty,email"`
	Name       *string  `json:"name,omitempty"`
	Preference *string  `json:"preference,omitempty"`
	WeightUnit *string  `json:"weightUnit,omitempty"`
	HeightUnit *string  `json:"heightUnit,omitempty"`
	Weight     *float64 `json:"weight,omitempty"`
	Height     *float64 `json:"height,omitempty"`
	ImageURI   *string  `json:"imageUri,omitempty"`
}

// Validate validates the UpdateUserRequest
func (r *UpdateUserRequest) Validate() error {
	if r.Preference != nil {
		if *r.Preference != PreferenceCardio && *r.Preference != PreferenceWeight {
			return errors.New("preference must be CARDIO or WEIGHT")
		}
	}

	if r.WeightUnit != nil {
		if *r.WeightUnit != WeightUnitKG && *r.WeightUnit != WeightUnitLBS {
			return errors.New("weightUnit must be KG or LBS")
		}
	}

	if r.HeightUnit != nil {
		if *r.HeightUnit != HeightUnitCM && *r.HeightUnit != HeightUnitINCH {
			return errors.New("heightUnit must be CM or INCH")
		}
	}

	return nil
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID         uint     `json:"id"`
	Email      string   `json:"email"`
	Name       *string  `json:"name"`
	Preference *string  `json:"preference"`
	WeightUnit *string  `json:"weightUnit"`
	HeightUnit *string  `json:"heightUnit"`
	Weight     *float64 `json:"weight"`
	Height     *float64 `json:"height"`
	ImageURI   *string  `json:"imageUri"`
}
