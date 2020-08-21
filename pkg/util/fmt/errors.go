package fmt

import (
	"fmt"
	"strings"
)

// FormatErrors returns a formatted list of error(s).
func FormatErrors(errs []error) string {
	if len(errs) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n\n", errs[0])
	}
	errStrs := make([]string, len(errs))
	for i, err := range errs {
		errStrs[i] = fmt.Sprintf("* %s", err)
	}
	return fmt.Sprintf("%d errors occurred:\n\t%s\n\n", len(errs), strings.Join(errStrs, "\n\t"))
}
