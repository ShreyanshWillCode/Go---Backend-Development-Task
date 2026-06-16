package models

import (
	"errors"
	"time"
)

type CreateUserRequest struct {
	Name	string	`json:"name" validate:"required,min=1,max=255"`

	Dob	string	`json:"dob" validate:"required,datetime=2006-01-02"`
}

type UpdateUserRequest struct {
	Name	string	`json:"name" validate:"required,min=1,max=255"`

	Dob	string	`json:"dob" validate:"required,datetime=2006-01-02"`
}

type UserResponse struct {
	ID	int32	`json:"id"`
	Name	string	`json:"name"`

	Dob	string	`json:"dob"`
}

type UserWithAgeResponse struct {
	ID	int32	`json:"id"`
	Name	string	`json:"name"`
	Dob	string	`json:"dob"`

	Age	int	`json:"age"`
}

type PaginatedUsersResponse struct {
	Data	[]UserWithAgeResponse	`json:"data"`
	Total	int64			`json:"total"`
	Page	int			`json:"page"`
	Limit	int			`json:"limit"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// ── Shared Helpers ─────────────────────────────────────────────────────────────

var ErrFutureDob = errors.New("date of birth cannot be in the future")

const (
	// DateLayout is the canonical date format used throughout the API.
	// It matches Go's reference time layout for "YYYY-MM-DD".
	DateLayout = "2006-01-02"
)

// ParseDob parses a "YYYY-MM-DD" string into a time.Time value (UTC midnight).
func ParseDob(raw string) (time.Time, error) {
	t, err := time.Parse(DateLayout, raw)
	if err != nil {
		return time.Time{}, err
	}
	if t.After(time.Now()) {
		return time.Time{}, ErrFutureDob
	}
	return t, nil
}
