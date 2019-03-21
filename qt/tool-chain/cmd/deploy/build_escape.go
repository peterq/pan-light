package deploy

import (
	"fmt"
	"strings"
)

func escapeFlags(ldFlags []string, ldFlagsCustom string) string {
	for _, s := range []string{"\"", "'"} {
		var newldFlagsCustom []string
		var insideQuotes bool
		for i, f := range strings.Split(ldFlagsCustom, s) {
			if i > 0 {
				if !insideQuotes {
					if strings.Contains(f, " ") {
						insideQuotes = true
						f = strings.Replace(f, " ", "_DONT_ESCAPE_", -1)
					}
				} else {
					insideQuotes = false
				}
			}
			newldFlagsCustom = append(newldFlagsCustom, f)
		}
		if len(newldFlagsCustom) > 0 {
			ldFlagsCustom = strings.Join(newldFlagsCustom, "")
		}
	}

	if len(ldFlagsCustom) > 0 {
		ldFlags = append(ldFlags, strings.Split(ldFlagsCustom, " ")...)
	}

	if out := strings.Replace(strings.Join(ldFlags, "\" \""), "_DONT_ESCAPE_", " ", -1); len(out) > 0 {
		return fmt.Sprintf("\"%v\"", out)
	}
	return ""
}
