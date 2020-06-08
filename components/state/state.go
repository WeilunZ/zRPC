package state

import "fmt"

const (
	OK             = 0
	FrameworkError = 1
	BusinuessError = 2
	InternalError  = 3
)

const (
	SUCCESS              = "success"
	InternalErrorMessage = "server internal error"
)

// Error defines all errors in the framework
type Error struct {
	Code    uint32
	Type    int
	Message string
}

func (e *Error) Error() string {
	if e == nil {
		return SUCCESS
	}
	if e.Type == FrameworkError {
		return fmt.Sprintf("type : framework, code : %d, msg : %s", e.Code, e.Message)
	}
	return fmt.Sprintf("type : business, code : %d, msg : %s", e.Code, e.Message)
}

// new a framework type error
func NewFrameworkError(code uint32, msg string) *Error {
	return &Error{
		Type:    FrameworkError,
		Code:    code,
		Message: msg,
	}
}

// new a business type error
func New(code uint32, msg string) *Error {
	return &Error{
		Type:    BusinuessError,
		Code:    code,
		Message: msg,
	}
}
