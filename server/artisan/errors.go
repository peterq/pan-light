package artisan

type AppError interface {
	error
	Code() int
	Parent() error
}

type NormalError struct {
	code    int
	message string
	parent  error
}

func (e *NormalError) Error() string {
	if e.parent != nil {
		return e.message + ": " + e.parent.Error()
	}
	return e.message
}

func (e *NormalError) Code() int {
	return e.code
}

func (e *NormalError) Parent() error {
	return e.parent
}

func NewError(message string, code int, parent error) *NormalError {
	return &NormalError{
		code:    code,
		message: message,
		parent:  parent,
	}
}

func ErrorFrom(err error) AppError {
	if e, ok := err.(AppError); ok {
		return e
	}
	return NewError(err.Error(), -1, nil)
}
