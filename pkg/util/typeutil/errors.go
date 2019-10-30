package typeutil

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type InvalidVersionKindError struct {
	ExpectedVersion, ExpectedKind string
	GotVersion, GotKind           string
}

func (e *InvalidVersionKindError) Error() string {
	m := fmt.Sprintf("typeutil: expected version %q with kind %q", e.ExpectedVersion, e.ExpectedKind)
	if e.GotVersion != "" || e.GotKind != "" {
		m += fmt.Sprintf(", but got version %q with kind %q", e.GotVersion, e.GotKind)
	}

	return m
}

type FieldValidationError struct {
	Context     string
	Field       string
	Description string
}

func (e *FieldValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Description)
}

type ValidationError struct {
	FieldErrors []*FieldValidationError
}

func (e *ValidationError) Error() string {
	if len(e.FieldErrors) == 0 {
		return "typeutil: general validation error"
	}

	fes := make([]string, len(e.FieldErrors))
	for i, fe := range e.FieldErrors {
		fes[i] = fe.Error()
	}

	return fmt.Sprintf("typeutil: validation error:\n* %s", strings.Join(fes, "\n* "))
}

func ValidationErrorFromResult(result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}

	errs := result.Errors()

	fes := make([]*FieldValidationError, len(errs))
	for i, err := range errs {
		fes[i] = &FieldValidationError{
			Context:     err.Context().String(),
			Field:       err.Field(),
			Description: err.Description(),
		}
	}

	return &ValidationError{
		FieldErrors: fes,
	}
}
