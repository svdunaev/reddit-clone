package domain

import "errors"

type PostValidationErrorField = string
type PostValidationErrorReason = string

type ValidationError struct {
	Field  PostValidationErrorField
	Reason PostValidationErrorReason
}

var ErrValidation = errors.New("validation failed")
