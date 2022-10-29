package pl

import "fmt"

type FuncMap map[string]any

func NewFuncMap() FuncMap {
	return FuncMap{
		"pass":   fnPass,
		"printf": fnPrintf,
	}
}

func fnPass(vs ...any) []any {
	return vs
}

func fnPrintf(format string, vs ...any) string {
	return fmt.Sprintf(format, vs...)
}
