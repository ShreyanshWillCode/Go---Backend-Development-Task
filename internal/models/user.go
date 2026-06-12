package models

import "time"

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

const (
	DateLayout = "2006-01-02"
)

func ParseDob(raw string) (time.Time, error) {
	return time.Parse(DateLayout, raw)
}
