package constants

import "errors"

var (
	ErrUsernameGuest = errors.New(`username cannot be "guest"`)
)
