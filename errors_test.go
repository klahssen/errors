package errors

import (
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestOrigin(t *testing.T) {
	tests := []struct {
		e      *Err
		origin error
	}{
		{
			e:      nil,
			origin: nil,
		},
		{
			e:      New(TypeInternal, "delete", nil).(*Err),
			origin: nil,
		},
		{
			e:      New(TypeInternal, "delete", New(TypeNotFound, "db", fmt.Errorf("failed"))).(*Err),
			origin: fmt.Errorf("failed"),
		},
		{
			e:      New(TypeInternal, "delete", fmt.Errorf("failed")).(*Err),
			origin: fmt.Errorf("failed"),
		},
	}
	for ind, test := range tests {
		o := test.e.Origin()
		if !reflect.DeepEqual(o, test.origin) {
			t.Errorf("test %d: expected origin %v received %v", ind, test.origin, o)
		}
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		err error
		t   ErrType
		ok  bool
	}{
		{
			err: nil,
			t:   TypeNotFound,
			ok:  false,
		},
		{
			err: fmt.Errorf("a special error"),
			t:   TypeNotFound,
			ok:  false,
		},
		{
			err: New(TypeInternal, "delete", fmt.Errorf("internal failure")),
			t:   TypeNotFound,
			ok:  false,
		},
		{
			err: New(TypeNotFound, "get", fmt.Errorf("item not found")),
			t:   TypeNotFound,
			ok:  true,
		},
		{
			err: New(TypeOther, "delete", New(TypeNotFound, "get", fmt.Errorf("item not found"))),
			t:   TypeNotFound,
			ok:  true,
		},
	}
	for ind, test := range tests {
		ok := Is(test.t, test.err)
		if ok != test.ok {
			t.Errorf("test %d: expected ok %v received %v", ind, test.ok, ok)
		}
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		e      *Err
		expect string
	}{
		{
			e:      nil,
			expect: "no error",
		},
		{
			e:      New(TypeNotFound, "get", nil).(*Err),
			expect: fmt.Sprintf("get: %s", TypeNotFound.String()),
		},
		{
			e:      New(TypeNotFound, "get", New(TypeInternal, "fetch", nil)).(*Err),
			expect: fmt.Sprintf("get: %s => fetch: internal error", TypeNotFound.String()),
		},
		{
			e:      New(TypeNotFound, "get", New(TypeInternal, "fetch", fmt.Errorf("connection refused"))).(*Err),
			expect: fmt.Sprintf("get: %s => fetch: internal error: connection refused", TypeNotFound.String()),
		},
	}

	for ind, test := range tests {
		msg := test.e.Error()
		if msg != test.expect {
			t.Errorf("test %d: expected msg '%s' received '%s'", ind, test.expect, msg)
		}
	}
}

func TestGetGRPCCode(t *testing.T) {
	tests := []struct {
		err  error
		code codes.Code
	}{
		{err: nil, code: codes.OK},
		{err: New(TypeOther, "", nil), code: codes.Internal},
		{err: New(TypeInternal, "", nil), code: codes.Internal},
		{err: New(TypeInvalidArg, "", nil), code: codes.InvalidArgument},
		{err: New(TypeInvalidOp, "", nil), code: codes.NotFound},
		{err: New(TypeNotFound, "", nil), code: codes.NotFound},
		{err: New(TypeAlreadyExists, "", nil), code: codes.AlreadyExists},
		{err: New(TypePermission, "", nil), code: codes.PermissionDenied},
		{err: New(TypeIO, "", nil), code: codes.Internal},
		{err: New(TypeTimeout, "", nil), code: codes.DeadlineExceeded},
		{err: New(TypeTooMany, "", nil), code: codes.ResourceExhausted},
		{err: New(TypeUnexpected, "", nil), code: codes.Internal},
		{err: New(TypeUnauthenticated, "", nil), code: codes.Unauthenticated},
	}
	var code codes.Code
	for ind, test := range tests {
		code = GetGRPCCode(test.err)
		if test.code != code {
			t.Errorf("test %d: expected Code %d:%s received %d:%s", ind, test.code, test.code.String(), code, code.String())
		}
	}
}
