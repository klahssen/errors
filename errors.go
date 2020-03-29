package errors

import (
	"bytes"

	"google.golang.org/grpc/codes"
)

const separator = " => " // ":\n\t"

// ErrType is the error type
type ErrType uint8

// error types
const (
	TypeOther      ErrType = iota // default
	TypeInternal                  // internal error
	TypeInvalidArg                // invalid arguments for an operation
	TypeInvalidRequestBody
	TypeInvalidOp     // invalid operation (method not allowed for exampled)
	TypeNotFound      // resource not found
	TypeAlreadyExists // conflict with existing resource
	TypePermission    // not authorized or unauthenticated
	TypeIO            // external io error (network failure etc)
	TypeTimeout
	TypeTooMany         // overload
	TypeUnexpected      // should be escalated
	TypeUnauthenticated // unauthenticated
)

// String returns the printable version of the error type
func (t ErrType) String() string {
	switch t {
	case TypeInternal:
		return "internal error"
	case TypeInvalidOp:
		return "invalid"
	case TypeInvalidArg:
		return "invalid argument(s)"
	case TypeInvalidRequestBody:
		return "invalid request body"
	case TypeUnexpected:
		return "unexpected error"
	case TypeNotFound:
		return "resource not found"
	case TypeAlreadyExists:
		return "conflict with existing resource"
	case TypePermission:
		return "permission denied"
	case TypeUnauthenticated:
		return "unauthenticated"
	case TypeTimeout:
		return "request timeout"
	case TypeTooMany:
		return "overload"
	case TypeIO:
		return "io error"
	case TypeOther:
		return "error"
	default:
		return "unknown error type"
	}
}

var (
	_ error = (*Err)(nil) // check if the Err type satisfies error
)

// Err holds an error
type Err struct {
	Typ   ErrType // error type
	Op    string  // operation that failed
	Cause error   // cause
}

// GetErrType returns the ErrType. Returns TypeOther if not castable as Err
func GetErrType(err error) ErrType {
	e, ok := err.(*Err)
	if !ok {
		return TypeOther
	}
	return e.Typ
}

// Is checks if err is of error type t
func Is(t ErrType, err error) bool {
	e, ok := err.(*Err)
	if !ok {
		return false
	}
	if e.Typ != TypeOther {
		return e.Typ == t
	}
	if e.Cause != nil {
		return Is(t, e.Cause)
	}
	return false
}

// New constructs a new error
func New(t ErrType, op string, cause error) error {
	return &Err{Typ: t, Cause: cause, Op: op}
}

// FirstError create a neew error with the first non nil error
func FirstError(t ErrType, op string, causes ...error) error {
	for _, cause := range causes {
		if cause != nil {
			return &Err{Typ: t, Cause: cause, Op: op}
		}
	}
	return nil
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func (e *Err) isZero() bool {
	return e.Op == "" && e.Typ == 0 && e.Cause == nil
}

// Origin returns the deepest error unwrapped from inner Cause error recursively
func (e *Err) Origin() error {
	if e == nil {
		return nil
	}
	if e.Cause == nil {
		return nil
	}
	if err, ok := e.Cause.(*Err); ok {
		return err.Origin()
	}
	return e.Cause
}

// Error to implement error interface. Tries to write <operation>: <error_type>: cause
func (e *Err) Error() string {
	if e == nil {
		return "no error"
	}
	b := new(bytes.Buffer)
	if e.Op != "" {
		pad(b, ": ")
		b.WriteString(e.Op)
	}
	if e.Typ > 0 {
		pad(b, ": ")
		b.WriteString(e.Typ.String())
	}
	if e.Cause != nil {
		if c, ok := e.Cause.(*Err); ok {
			if !c.isZero() {
				pad(b, separator)
				b.WriteString(c.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.Cause.Error())
		}
	}

	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

//GetGRPCCode tries to cast error as *Err and return a gRPC code based on internal ErrType
func GetGRPCCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	e, ok := err.(*Err)
	if !ok {
		return codes.Internal
	}
	if is(TypeInvalidArg, e) {
		return codes.InvalidArgument
	}
	if is(TypeNotFound, e) {
		return codes.NotFound
	}
	if is(TypeInvalidOp, e) {
		return codes.NotFound
	}
	if is(TypeTimeout, e) {
		return codes.DeadlineExceeded
	}
	if is(TypePermission, e) {
		return codes.PermissionDenied
	}
	if is(TypeUnauthenticated, e) {
		return codes.Unauthenticated
	}
	if is(TypeTooMany, e) {
		return codes.ResourceExhausted
	}
	if is(TypeAlreadyExists, e) {
		return codes.AlreadyExists
	}

	return codes.Internal
}

func is(t ErrType, e *Err) bool {
	if e.Typ != TypeOther {
		return e.Typ == t
	}
	if e.Cause != nil {
		return Is(t, e.Cause)
	}
	return false
}
