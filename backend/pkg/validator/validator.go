// Package validator provides functions for validating data.
package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(s any) error
}

type CustomValidator struct {
	v *validator.Validate
}

var (
	once     sync.Once
	instance *CustomValidator
)

func New() Validator {
	once.Do(func() {
		v := validator.New()
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		instance = &CustomValidator{v: v}
	})
	return instance
}

type Error struct {
	Messages []string
}

func (e *Error) Error() string {
	return strings.Join(e.Messages, " | ")
}

func (cv *CustomValidator) Struct(s any) error {
	err := cv.v.Struct(s)
	if err == nil {
		return nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range errs {
			msg := fmt.Sprintf("%s: %s", e.Field(), formatMessage(e))
			messages = append(messages, msg)
		}

		return &Error{Messages: messages}
	}

	return err
}

func formatMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "required field"
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("minimum length is %s", e.Param())
	case "max":
		return fmt.Sprintf("maximum length is %s", e.Param())
	default:
		return fmt.Sprintf("invalid value (%s)", e.Tag())
	}
}
