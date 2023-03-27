package errors

// Custom Application level error interface
type AppErrorInterface interface {
	error
	GetError() error
	GetMessage() string
}

// Base Aplication level error
type ErrBase struct {
	Err     error  `json:"-"`
	Message string `json:"message"`
}

func (e ErrBase) GetMessage() string {
	return e.Message
}

func (e ErrBase) GetError() error {
	return e.Err
}

func (e ErrBase) Error() string {
	return e.GetMessage()
}

// Custom Application level Not Found error
type ErrResourceNotFound struct{ ErrBase }

func NewErrResourceNotFound(err error, msg string) error {
	return ErrResourceNotFound{ErrBase: ErrBase{Err: err, Message: msg}}
}

// Custom Application level Duplicate Entry error
type ErrDuplicateEntry struct{ ErrBase }

func NewErrDuplicateEntry(err error, msg string) error {
	return ErrDuplicateEntry{ErrBase: ErrBase{Err: err, Message: msg}}
}

// Custom Application level Unauthorized error
type ErrUnauthorized struct{ ErrBase }

func NewErrUnauthorized(err error, msg string) error {
	return ErrUnauthorized{ErrBase: ErrBase{Err: err, Message: msg}}
}

// Custom Application level Forbidden error
type ErrForbidden struct{ ErrBase }

func NewErrForbidden(err error, msg string) error {
	return ErrForbidden{ErrBase: ErrBase{Err: err, Message: msg}}
}

// Custom Application level BadRequest error
type ErrBadRequest struct{ ErrBase }

func NewErrBadRequest(err error, msg string) error {
	return ErrBadRequest{ErrBase: ErrBase{Err: err, Message: msg}}
}
