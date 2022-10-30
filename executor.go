package pl

import (
	"errors"
	"fmt"
	"reflect"
)

type fnNode struct {
	name string
	args []any

	has_nested bool
}

type Executor struct {
	Funcs FuncMap
	Props Props
	Convs ConvMap
}

func NewExecutor() *Executor {
	return &Executor{
		Funcs: NewFuncMap(),
		Props: NewProps(M{}),
		Convs: NewConvMap(),
	}
}

func (e *Executor) ExecuteExpr(expr string) ([]any, error) {
	pl, err := ParseString(expr)
	if err != nil {
		return nil, err
	}

	return e.Execute(pl)
}

func (e *Executor) Execute(pl *Pl) ([]any, error) {
	args_prev := []any{}
	for i, fn := range pl.Funcs {
		err := func() error {
			f, ok := e.Funcs[fn.Name]
			if !ok {
				return errors.New("not defined")
			}

			fnode, err := e.evaluateFn(fn)
			if err != nil {
				return err
			}

			args := make([]any, 0, len(fnode.args)+len(args_prev))
			if !fnode.has_nested {
				args = append(args, fnode.args...)
			} else {
				for i, arg := range fnode.args {
					nested, ok := arg.(*Pl)
					if !ok {
						args = append(args, arg)
						continue
					}

					rst, err := e.Execute(nested)
					if err != nil {
						return fmt.Errorf("arg[%d]: %w", i, err)
					}

					args = append(args, rst...)
				}
			}

			args = append(args, args_prev...)

			rst, err := e.invokeFn(f, args)
			if err != nil {
				return err
			}

			if reflect.TypeOf(rst).Kind() != reflect.Slice {
				args_prev = []any{rst}
			} else {
				v := reflect.ValueOf(rst)
				args_prev = make([]any, v.Len())
				for i := 0; i < v.Len(); i++ {
					args_prev[i] = v.Index(i).Interface()
				}
			}

			return nil
		}()

		if err != nil {
			return nil, fmt.Errorf("fn[%d] %s: %w", i, fn.Name, err)
		}
	}

	return args_prev, nil
}

func (e *Executor) evaluateFn(fn *Fn) (*fnNode, error) {
	rst := &fnNode{name: fn.Name, args: make([]any, len(fn.Args))}
	for i, arg := range fn.Args {
		if arg.String != nil {
			rst.args[i] = *arg.String
		} else if arg.Float != nil {
			rst.args[i] = *arg.Float
		} else if arg.Int != nil {
			rst.args[i] = *arg.Int
		} else if arg.Ref != nil {
			arg, err := e.Props.Resolve(arg.Ref)
			if err != nil {
				return nil, fmt.Errorf("arg[%d]: reference: %w", i, err)
			}

			rst.args[i] = arg
		} else if arg.Nested != nil {
			rst.args[i] = arg.Nested
			rst.has_nested = true
		} else {
			return nil, fmt.Errorf("arg[%d]: empty value", i)
		}
	}

	return rst, nil
}

func (e *Executor) invokeFn(fn any, args []any) (any, error) {
	fv := reflect.ValueOf(fn)
	ft := fv.Type()

	// Check if the number of returned values is valid.
	if n := ft.NumOut(); n > 2 || n == 0 {
		return nil, fmt.Errorf("function have to return one or two values but %d values are returned", n)
	} else if n == 2 && !ft.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, fmt.Errorf("type of second return value of the function must be an error but it was %s", ft.Out(1).Name())
	}

	// Check if the number of argument is fit.
	num_fixed_args := ft.NumIn()
	if ft.IsVariadic() {
		num_fixed_args--
		if len(args) < num_fixed_args {
			return nil, fmt.Errorf("expected at least %d args but %d args are given", num_fixed_args, len(args))
		}
	} else if len(args) != num_fixed_args {
		return nil, fmt.Errorf("expected %d args but %d args are given", len(args), num_fixed_args)
	}

	input_args := make([]reflect.Value, len(args))
	for i, arg := range args {
		j := i
		if i >= num_fixed_args {
			j = num_fixed_args
		}

		t_arg := reflect.TypeOf(arg)
		t_in := ft.In(j)
		if i >= num_fixed_args {
			t_in = t_in.Elem()
		}
		if t_arg.AssignableTo(t_in) {
			input_args[i] = reflect.ValueOf(arg)
			continue
		}

		v_arg := reflect.ValueOf(arg)
		v, err := e.Convs.Convert(t_in, t_arg, v_arg)
		if err != nil {
			// Try to call String() if parameter type is string.
			s, ok := func() (string, bool) {
				if t_in.Kind() != reflect.String {
					return "", false
				}

				if s, ok := v_arg.Interface().(interface{ String() string }); ok {
					return s.String(), true
				}

				p := reflect.New(t_arg)
				p.Elem().Set(v_arg)

				if s, ok := p.Interface().(interface{ String() string }); ok {
					return s.String(), true
				} else {
					return "", false
				}
			}()
			if !ok {
				return nil, fmt.Errorf("arg[%d]: %w", i, err)
			}

			v = s
		}

		input_args[i] = reflect.ValueOf(v)
	}

	rst := fv.Call(input_args)
	if len(rst) == 1 || (len(rst) == 2 && rst[1].IsNil()) {
		return rst[0].Interface(), nil
	} else {
		err := rst[1].Interface().(error)
		return rst[0].Interface(), err
	}
}
