package v1

import (
	"fmt"
	"strings"
)

type InvalidVersionKindError struct {
	ExpectedVersion, ExpectedKind string
	GotVersion, GotKind           string
}

func (e *InvalidVersionKindError) Error() string {
	m := fmt.Sprintf("v1: expected version %q with kind %q", e.ExpectedVersion, e.ExpectedKind)
	if e.GotVersion != "" || e.GotKind != "" {
		m += fmt.Sprintf(", but got version %q with kind %q", e.GotVersion, e.GotKind)
	}

	return m
}

type SchemaValidationFieldError struct {
	Context     string
	Field       string
	Description string
}

func (e *SchemaValidationFieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Description)
}

type SchemaValidationError struct {
	FieldErrors []*SchemaValidationFieldError
}

func (e *SchemaValidationError) Error() string {
	if len(e.FieldErrors) == 0 {
		return "v1: general schema validation error"
	}

	fes := make([]string, len(e.FieldErrors))
	for i, fe := range e.FieldErrors {
		fes[i] = fe.Error()
	}

	return fmt.Sprintf("v1: schema validation error:\n* %s", strings.Join(fes, "\n* "))
}
