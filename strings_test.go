package errors

import "testing"

func TestStringError(t *testing.T) {
	tests := []struct {
		s      String
		expect string
	}{
		{
			s:      " hello ",
			expect: " hello ",
		},
		{
			s:      "",
			expect: "",
		},
	}
	var s string
	for ind, test := range tests {
		s = test.s.Error()
		if s != test.expect {
			t.Errorf("test %d: expected error '%s' received '%s'", ind, test.expect, s)
		}
	}
}

func TestStringStringer(t *testing.T) {
	tests := []struct {
		s      String
		expect string
	}{
		{
			s:      " hello ",
			expect: " hello ",
		},
		{
			s:      "",
			expect: "",
		},
	}
	var s string
	for ind, test := range tests {
		s = test.s.String()
		if s != test.expect {
			t.Errorf("test %d: expected string '%s' received '%s'", ind, test.expect, s)
		}
	}
}
