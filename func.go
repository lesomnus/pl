package pl

import (
	"github.com/lesomnus/pl/funcs"
)

type FuncMap map[string]any

func NewFuncMap() FuncMap {
	return FuncMap{
		"pass":   funcs.Pass,
		"printf": funcs.Printf,
		"regex":  funcs.Regex,
	}
}
