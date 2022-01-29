package forms

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Form struct {
	Errors errors
	url.Values
}

func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(errors),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) NoWhiteSpace(field string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, c := range value {
		if unicode.IsSpace(c) {
			f.Errors.Add(field, "This field cannot contain whitespace")
			return
		}
	}
}

func (f *Form) MinLength(field string, minLength int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < minLength {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum length is %d)", minLength))
	}
}

func (f *Form) RequireTypeInt(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if value == "" {
			return
		}

		_, err := strconv.Atoi(value)
		if err != nil {
			f.Errors.Add(field, "Please enter a valid number")
		}
	}
}

func (f *Form) MaxLength(field string, maxLength int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > maxLength {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum length is %d)", maxLength))
	}
}

func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}
