package chopchop

type Error struct {
	Code int
	Message string
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Message: msg,
	}
}
