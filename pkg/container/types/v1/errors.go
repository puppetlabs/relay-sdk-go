package v1

import "fmt"

type InvalidVersionKindError struct {
	ExpectedVersion, ExpectedKind string
	GotVersion, GotKind           string
}

func (e *InvalidVersionKindError) Error() string {
	m := fmt.Sprintf("expected version %q with kind %q", e.ExpectedVersion, e.ExpectedKind)
	if e.GotVersion != "" || e.GotKind != "" {
		m += fmt.Sprintf(", but got version %q with kind %q", e.GotVersion, e.GotKind)
	}

	return m
}
