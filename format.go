package multierr

import (
	"fmt"
	"strings"
)

// FormatterFunc is called by Error.Error() to convert
// multi-errors into a human readable strings.
type FormatterFunc func([]error) string

// ListFormatterFunc puts each sub-error in a new line.
// All errors will be indented and titled with a generic "n errors occurred".
func ListFormatterFunc(errs []error) string {
	if len(errs) == 0 {
		return "no errors occurred"
	}
	//if len(errs) == 1 { // This might yield better results in most cases, but it would be breaking and surprising - it would make things more difficult to understand and reason about.
	//	return errs[0].Error()
	//}

	plural := "errors"
	if len(errs) == 1 {
		plural = "error"
	}
	title := fmt.Sprintf("%d %s occurred:", len(errs), plural)
	return TitledListFormatter(title)(errs)
}

// TitledListFormatter returns a formatter func that puts each sub-error in a new, indented line.
// The errors are titled with the given text.
func TitledListFormatter(title string) FormatterFunc {
	return func(errs []error) string {
		if len(errs) == 0 {
			return "no errors occurred"
		}

		var str = title
		for _, err := range errs {
			msg := strings.Replace(err.Error(), "\n", "\n    ", -1)
			str += "\n  - " + msg
		}
		return str
	}
}

// PrefixedListFormatter returns a formatter func that puts each sub-error in a new line.
// The errors are prefixed with the given text.
func PrefixedListFormatter(prefix string) FormatterFunc {
	return func(errs []error) string {
		if len(errs) == 0 {
			return "no errors occurred"
		}
		multilineIndent := strings.Repeat(" ", len(prefix))

		var str string
		for _, err := range errs {
			msg := strings.Replace(err.Error(), "\n", "\n"+multilineIndent, -1)
			str += "\n" + prefix + msg
		}
		return str[1:]
	}
}
