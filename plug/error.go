package plug

import (
	"fmt"
	"strings"
)

// ExitError .
type ExitError struct {
	Text     string
	ExitCode int
}

func (e ExitError) Error() string {
	return e.Text
}

var ErrUsageError = ExitError{
	Text:     "usage error",
	ExitCode: 1,
}

// ExecError .
type ExecError struct {
	Err         error
	UsageErrors map[string][]string
}

func (e ExecError) Error() string {
	var errs []string
	if len(e.UsageErrors) > 0 {
		errs = append(errs, fmt.Sprintf("%v usage errors", len(e.UsageErrors)))
		return "usage errors"
	}
	if e.Err != nil {
		errs = append(errs, e.Err.Error())
	}

	if len(errs) == 0 {
		return "no errors"
	}

	return strings.Join(errs, "; ")
}
