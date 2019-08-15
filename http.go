package errors

import (
	"net/http"
)

//HTTPErr is a http error
type HTTPErr struct {
	status int
	err    error
}

//Error to satisfy the error interface
func (h *HTTPErr) Error() string {
	if h == nil || h.err == nil {
		return ""
	}
	return h.err.Error()
}

//Status returns the inner http status code
func (h *HTTPErr) Status() int {
	return h.status
}

//Cause returns the inner error
func (h *HTTPErr) Cause() error {
	return h.err
}

//NewHTTPErr returns a new HttpErr with embedded err and http status code. If invalid status code => default to http.StatusBadRequest
func NewHTTPErr(status int, err error) error {
	if err == nil {
		return nil
	}
	if http.StatusText(status) == "" {
		//invalid status code
		status = http.StatusBadRequest
	}
	return &HTTPErr{
		status: status,
		err:    err,
	}
}

//Cause will return the deepest error
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
