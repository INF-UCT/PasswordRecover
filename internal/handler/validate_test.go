package handler

import (
	"reflect"
	"sort"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name      string
		pw        string
		wantCodes []string
	}{
		{
			name:      "empty",
			pw:        "",
			wantCodes: []string{"too_short", "missing_lowercase", "missing_uppercase", "missing_digit", "missing_symbol"},
		},
		{
			name:      "single_char",
			pw:        "a",
			wantCodes: []string{"too_short", "missing_uppercase", "missing_digit", "missing_symbol"},
		},
		{
			name:      "too_short_only_lower",
			pw:        "abc",
			wantCodes: []string{"too_short", "missing_uppercase", "missing_digit", "missing_symbol"},
		},
		{
			name:      "no_uppercase",
			pw:        "abcdefg1!",
			wantCodes: []string{"missing_uppercase"},
		},
		{
			name:      "no_lowercase",
			pw:        "ABCDEFG1!",
			wantCodes: []string{"missing_lowercase"},
		},
		{
			name:      "no_digit",
			pw:        "Abcdefgh!",
			wantCodes: []string{"missing_digit"},
		},
		{
			name:      "no_symbol",
			pw:        "Abcdefg1",
			wantCodes: []string{"missing_symbol"},
		},
		{
			name:      "no_upper_no_digit",
			pw:        "abcdefgh!",
			wantCodes: []string{"missing_uppercase", "missing_digit"},
		},
		{
			name:      "no_lower_no_symbol",
			pw:        "ABCDEFG1",
			wantCodes: []string{"missing_lowercase", "missing_symbol"},
		},
		{
			name:      "valid_ascii",
			pw:        "ValidPass1!",
			wantCodes: nil,
		},
		{
			name:      "valid_with_unicode",
			pw:        "MiClaveÑoño2024!",
			wantCodes: nil,
		},
		{
			name:      "valid_only_unicode_lower",
			pw:        "ñandú2024!A",
			wantCodes: nil,
		},
		{
			name:      "exactly_8_with_all_classes",
			pw:        "aA1!aA1!",
			wantCodes: nil,
		},
		{
			name:      "seven_chars_with_all_classes",
			pw:        "aA1!aA1",
			wantCodes: []string{"too_short"},
		},
		{
			name:      "only_specials",
			pw:        "!!!!!!!!",
			wantCodes: []string{"missing_lowercase", "missing_uppercase", "missing_digit"},
		},
		{
			name:      "space_is_not_special",
			pw:        "Abc defg1",
			wantCodes: []string{"missing_symbol"},
		},
		{
			name:      "tab_is_not_special",
			pw:        "Abc\tdefg1",
			wantCodes: []string{"missing_symbol"},
		},
		{
			name:      "double_quote_is_special",
			pw:        `Abc"defg1`,
			wantCodes: nil,
		},
		{
			name:      "backslash_is_special",
			pw:        `Abc\defg1`,
			wantCodes: nil,
		},
		{
			name:      "backtick_is_special",
			pw:        "Abc`defg1",
			wantCodes: nil,
		},
		{
			name:      "tilde_is_special",
			pw:        "Abc~defg1",
			wantCodes: nil,
		},
		{
			name:      "tilde_n_is_not_lowercase_ascii",
			pw:        "Abc~defg1",
			wantCodes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validatePassword(tt.pw)
			gotCodes := codesFromFieldErrors(got)
			wantCodes := append([]string(nil), tt.wantCodes...)
			sort.Strings(wantCodes)
			if !reflect.DeepEqual(gotCodes, wantCodes) {
				t.Errorf("validatePassword(%q) got codes %v, want %v (full errors: %+v)", tt.pw, gotCodes, wantCodes, got)
			}
			for _, e := range got {
				if e.Field != "new_password" {
					t.Errorf("validatePassword(%q): expected field=new_password, got %q", tt.pw, e.Field)
				}
				if e.Code == "" {
					t.Errorf("validatePassword(%q): empty code in error %+v", tt.pw, e)
				}
				if e.Message == "" {
					t.Errorf("validatePassword(%q): empty message in error %+v", tt.pw, e)
				}
			}
		})
	}
}

func TestValidatePasswordMatch(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		wantCode string
	}{
		{name: "equal", a: "Pass1!", b: "Pass1!", wantCode: ""},
		{name: "different", a: "Pass1!", b: "Pass2!", wantCode: "password_mismatch"},
		{name: "both_empty", a: "", b: "", wantCode: ""},
		{name: "case_sensitive", a: "pass1!", b: "Pass1!", wantCode: "password_mismatch"},
		{name: "unicode_equality", a: "MiClaveÑoño2024!", b: "MiClaveÑoño2024!", wantCode: ""},
		{name: "one_empty", a: "", b: "Pass1!", wantCode: "password_mismatch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validatePasswordMatch(tt.a, tt.b)
			if tt.wantCode == "" {
				if got != nil {
					t.Errorf("validatePasswordMatch(%q, %q): expected nil, got %+v", tt.a, tt.b, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("validatePasswordMatch(%q, %q): expected error with code %q, got nil", tt.a, tt.b, tt.wantCode)
			}
			if got.Code != tt.wantCode {
				t.Errorf("validatePasswordMatch(%q, %q): expected code %q, got %q", tt.a, tt.b, tt.wantCode, got.Code)
			}
			if got.Field != "confirmed_password" {
				t.Errorf("validatePasswordMatch(%q, %q): expected field=confirmed_password, got %q", tt.a, tt.b, got.Field)
			}
			if got.Message == "" {
				t.Errorf("validatePasswordMatch(%q, %q): empty message in error %+v", tt.a, tt.b, got)
			}
		})
	}
}

func TestUtf8Len(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{name: "empty", s: "", want: 0},
		{name: "ascii", s: "abcdef", want: 6},
		{name: "with_spaces", s: "ab cd", want: 5},
		{name: "with_unicode", s: "ñandú", want: 5},
		{name: "with_emoji", s: "ab🔒cd", want: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utf8Len(tt.s); got != tt.want {
				t.Errorf("utf8Len(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

func codesFromFieldErrors(errs []fieldError) []string {
	if len(errs) == 0 {
		return nil
	}
	codes := make([]string, 0, len(errs))
	for _, e := range errs {
		codes = append(codes, e.Code)
	}
	sort.Strings(codes)
	return codes
}
