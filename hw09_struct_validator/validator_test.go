package hw09structvalidator

import (
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	IDK struct {
		Text string `validate:"i-dont-know-this-rule:for-sure"`
	}

	BadFmt struct {
		SomeName string `validate:"im-wrong"`
	}

	LetsTestIntSlice struct {
		ImSlice []int `validate:"min:3|max:6"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "73fed80e-4bc0-11ee-be56-0242ac120002",
				Name:   "test",
				Age:    32,
				Email:  "mail@example.com",
				Role:   "admin",
				Phones: []string{"+7800655005"},
				meta:   []byte("raw raw raw"),
			},
			expectedErr: nil,
		},
		{
			in:          "test",
			expectedErr: ErrInvalidArg,
		},
		{
			in: App{Version: "123"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrValidateLen},
			},
		},
		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in:          IDK{Text: "text"},
			expectedErr: ErrUnknnownRule,
		},
		{
			in:          BadFmt{SomeName: "name"},
			expectedErr: ErrInvalidRule,
		},
		{
			in: User{
				ID:     "1234",
				Name:   "test",
				Age:    500,
				Email:  "mail-me-plz",
				Role:   "dude",
				Phones: []string{"655005"},
				meta:   []byte("raw raw raw"),
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrValidateLen},
				ValidationError{Field: "Age", Err: ErrValidateMax},
				ValidationError{Field: "Email", Err: ErrValidateRegex},
				ValidationError{Field: "Role", Err: ErrValidateIn},
				ValidationError{Field: "Phones[0]", Err: ErrValidateLen},
			},
		},
		{
			in:          Response{Code: 200, Body: ""},
			expectedErr: nil,
		},
		{
			in: Response{Code: 201, Body: ""},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrValidateIn},
			},
		},
		{
			in: User{
				ID:     "73fed80e-4bc0-11ee-be56-0242ac120002",
				Name:   "test",
				Age:    2,
				Email:  "mail@example.com",
				Role:   "admin",
				Phones: []string{"+7800655005"},
				meta:   []byte("raw raw raw"),
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrValidateMin},
			},
		},
		{
			in: LetsTestIntSlice{ImSlice: []int{1, 3, 6, 10}},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ImSlice[0]", Err: ErrValidateMin},
				ValidationError{Field: "ImSlice[3]", Err: ErrValidateMax},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			var valErrs ValidationErrors
			if errors.As(err, &valErrs) {
				var expErrs ValidationErrors
				ok := errors.As(tt.expectedErr, &expErrs)
				if !ok {
					t.Errorf("unexpected errors, got: %v, expected: %v", err, tt.expectedErr)
				}

				if len(expErrs) != len(valErrs) {
					t.Errorf("unexpected errors, got: %v, expected: %v", valErrs, expErrs)
				}

				for i, verr := range valErrs {
					if !(errors.Is(verr.Err, (&expErrs[i]).Err) && verr.Field == expErrs[i].Field) {
						t.Errorf("unexpected error, got: %v, expected: %v", verr, expErrs[i])
					}
				}
			} else if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error, got: %v, expected: %v", err, tt.expectedErr)
			}
		})
	}
}
