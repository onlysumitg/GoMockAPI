package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Add a new ErrInvalidCredentials error. We'll use this later if a user
	// tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("Invalid user credentials")
	// Add a new ErrDuplicateEmail error. We'll use this later if a user
	// tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
	ErrNotFound       = errors.New("models: Not found")

	ErrUserNotFound       = errors.New("User not found")
	ErrServerNotFound     = errors.New("Not Found")
	ErrSavedQueryNotFound = errors.New("models: Saved query not found")
)