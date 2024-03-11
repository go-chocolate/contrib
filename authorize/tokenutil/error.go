package tokenutil

type textError string

func (e textError) Error() string {
	return string(e)
}

const (
	ErrTokenInvalid = textError("token invalid")
	ErrTokenExpired = textError("token expired")
)
