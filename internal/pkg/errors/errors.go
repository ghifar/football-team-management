package errors

import (
	"fmt"
	"net/http"
)

// AppError represents application-specific errors
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Common application errors
var (
	ErrTeamNotFound = &AppError{
		Code:       "TEAM_NOT_FOUND",
		Message:    "Team not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrTeamAlreadyExists = &AppError{
		Code:       "TEAM_ALREADY_EXISTS",
		Message:    "Team already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrPlayerNotFound = &AppError{
		Code:       "PLAYER_NOT_FOUND",
		Message:    "Player not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrPlayerAlreadyExists = &AppError{
		Code:       "PLAYER_ALREADY_EXISTS",
		Message:    "Player already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrJerseyNumberTaken = &AppError{
		Code:       "JERSEY_NUMBER_TAKEN",
		Message:    "Jersey number already taken in this team",
		HTTPStatus: http.StatusConflict,
	}

	ErrMatchNotFound = &AppError{
		Code:       "MATCH_NOT_FOUND",
		Message:    "Match not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrMatchResultNotFound = &AppError{
		Code:       "MATCH_RESULT_NOT_FOUND",
		Message:    "Match result not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrMatchResultAlreadyExists = &AppError{
		Code:       "MATCH_RESULT_ALREADY_EXISTS",
		Message:    "Result already exists for this match",
		HTTPStatus: http.StatusConflict,
	}

	ErrInvalidInput = &AppError{
		Code:       "INVALID_INPUT",
		Message:    "Invalid input data",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized access",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "Insufficient permissions",
		HTTPStatus: http.StatusForbidden,
	}

	ErrInternalServer = &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// NewAppError creates a new application error
func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, message string) error {
	if appErr, ok := err.(*AppError); ok {
		appErr.Details = message
		return appErr
	}
	return &AppError{
		Code:       "WRAPPED_ERROR",
		Message:    message,
		Details:    err.Error(),
		HTTPStatus: http.StatusInternalServerError,
	}
}
