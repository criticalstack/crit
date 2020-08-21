package fmt

import "strconv"

func Unquote(s string) string {
	us, err := strconv.Unquote(s)
	if err != nil {
		return s
	}
	return us
}
