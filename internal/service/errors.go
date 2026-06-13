package service

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func validationError(message string) *Error {
	return &Error{
		Code:    "validation_error",
		Message: message,
	}
}
