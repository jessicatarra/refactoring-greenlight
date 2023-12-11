package domain

import "errors"

var (
	ErrEditConflict          = errors.New("edit conflict")
	ErrRecordNotFound        = errors.New("record not found")
	ErrDuplicateEmail        = errors.New("duplicate email")
	ErrPermissionNotIncluded = errors.New("permission not included")
)
