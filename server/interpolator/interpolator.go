package interpolator

import (
	"regexp"
)

type Interpolator struct {
	Re          *regexp.Regexp
	Matches     map[string]string
	ReplaceFunc func(s string) string
}

func (i *Interpolator) Interpolate(input string) string {
	return i.Re.ReplaceAllStringFunc(input, i.ReplaceFunc)
}
