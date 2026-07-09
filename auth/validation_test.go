package auth

import "testing"

func TestIsNotBlank(t *testing.T) {

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{
			name:  "valid username",
			value: "john",
			valid: true,
		},
		{
			name:  "empty username",
			value: "",
			valid: false,
		},
		{
			name:  "whitespace only",
			value: "   ",
			valid: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsNotBlank(test.value)
			if got != test.valid {
				t.Errorf("expected %v, got %v", test.valid, got)
			}
		})
	}

}

func TestServerUrlValidation(t *testing.T) {

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{
			name:  "valid http url",
			value: "http://test.server",
			valid: true,
		},
		{
			name:  "valid https url",
			value: "https://test.server",
			valid: true,
		},
		{
			name:  "invalid protocol",
			value: "ftp://test.server",
			valid: false,
		},
		{
			name:  "missing protocol",
			value: "test.server",
			valid: false,
		},
		{
			name:  "invalid format",
			value: "http://test  server",
			valid: false,
		},
		{
			name:  "empty string",
			value: "         ",
			valid: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsServerUrlValid(test.value)
			if got != test.valid {
				t.Errorf("expected %v, got %v", test.valid, got)
			}
		})
	}
}
