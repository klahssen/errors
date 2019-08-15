package errors

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestCause(t *testing.T) {
	tests := []struct {
		err      error
		expected error
	}{
		{
			err:      nil,
			expected: nil,
		},
		{
			err:      &HTTPErr{err: fmt.Errorf("some cause error")},
			expected: fmt.Errorf("some cause error"),
		},
		{
			err:      &HTTPErr{err: nil},
			expected: nil,
		},
		{
			err:      fmt.Errorf("simple error"),
			expected: fmt.Errorf("simple error"),
		},
	}
	for ind, test := range tests {
		err := Cause(test.err)
		if !reflect.DeepEqual(err, test.expected) {
			t.Errorf("test %d: expected err %+v received %+v", ind, test.expected, err)
		}
	}
}

func TestHTTPError(t *testing.T) {
	tests := []struct {
		err      *HTTPErr
		expected string
	}{
		{
			err:      nil,
			expected: "",
		},
		{
			err: &HTTPErr{
				status: http.StatusBadRequest,
				err:    fmt.Errorf("simple error"),
			},
			expected: "simple error",
		},
		{
			err: &HTTPErr{
				status: http.StatusBadRequest,
				err:    errors.Wrap(fmt.Errorf("cause"), "msg"),
			},
			expected: "msg: cause",
		},
	}
	for ind, test := range tests {
		err := test.err.Error()
		if err != test.expected {
			t.Errorf("test %d: expected err '%s' received '%s'", ind, test.expected, err)
		}
	}
}

func TestNewHTTPErr(t *testing.T) {
	tests := []struct {
		status   int
		err      error
		expected error
	}{
		{
			status: http.StatusBadRequest,
			err:    fmt.Errorf("cause"),
			expected: &HTTPErr{
				err:    fmt.Errorf("cause"),
				status: http.StatusBadRequest,
			},
		},
		{
			status: -5,
			err:    fmt.Errorf("cause"),
			expected: &HTTPErr{
				err:    fmt.Errorf("cause"),
				status: http.StatusBadRequest,
			},
		},
		{
			status:   http.StatusBadRequest,
			err:      nil,
			expected: nil,
		},
	}
	for ind, test := range tests {
		err := NewHTTPErr(test.status, test.err)
		if !reflect.DeepEqual(err, test.expected) {
			t.Errorf("test %d: expected err %+v received %+v", ind, test.expected, err)
		}
	}
}
