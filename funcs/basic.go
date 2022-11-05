package funcs

import "fmt"

func Pass(vs ...any) []any {
	return vs
}

func Printf(format string, vs ...any) string {
	return fmt.Sprintf(format, vs...)
}
