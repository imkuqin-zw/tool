package validation

import (
	"testing"
	"regexp"
	"fmt"
)

func TestParseRegex(t *testing.T)  {
	regex := regexp.MustCompile(`^regex\((^/)[^\(\)]+\)$`)
	te := `regex(/tre{fdsfsdffsdf.e3wt/, "no")`
	fmt.Println(regex.MatchString(te))
}
