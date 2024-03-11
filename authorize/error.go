package authorize

import "errors"

var (
	ErrInvalidUsername        = errors.New("invalid username")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrUserNotFound           = errors.New("user not found")
	ErrUnsupportedContentType = errors.New("unsupported content type")
)
