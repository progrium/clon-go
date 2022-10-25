package clon

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func Parse(args []string) (any, error) {
	var pargs []arg
	for _, a := range args {
		pargs = append(pargs, parseArg(a))
	}
	if len(pargs[0].path) == 0 {
		var v []any
		for _, a := range pargs {
			if len(a.path) != 0 {
				m := make(map[string]any)
				if err := applyArg(m, a); err != nil {
					panic(err)
				}
				v = append(v, m)
			} else {
				v = append(v, a.value)
			}

		}
		return v, nil
	}
	v := make(map[string]any)
	if pargs[0].tla {
		v["tla"] = []any{}
	}
	for _, a := range pargs {
		if err := applyArg(v, a); err != nil {
			panic(err)
		}
	}
	if pargs[0].tla {
		return v["tla"], nil
	}
	return v, nil
}

type arg struct {
	path   []string
	value  any
	append bool
	tla    bool // top level array
}

func parseArg(s string) (a arg) {
	kv := strings.Split(s, "=")
	path := kv[0]
	if strings.HasSuffix(path, ":") {
		path = path[:len(path)-1]
		a.value = parseRaw(kv[1])
	} else {
		a.value = strings.Trim(kv[1], "'")
	}
	if strings.HasSuffix(path, "[]") {
		a.append = true
		path = path[:len(path)-2]
	}
	path = strings.ReplaceAll(path, "]", "")
	a.path = strings.Split(path, "[")
	if a.path[0] == "" {
		a.path = append([]string{"tla"}, a.path[1:]...)
		a.tla = true
	}
	return a
}

func parseRaw(v string) any {
	if v == "true" {
		return true
	}
	if v == "false" {
		return false
	}
	if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
		v = strings.Trim(v, "'")
		var vv any
		if err := json.Unmarshal([]byte(v), &vv); err != nil {
			panic(err)
		}
		return vv
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return i
}

func applyArg(m map[string]any, v arg) error {
	var cur any
	cur = m
	for idx, part := range v.path {
		isLast := idx == len(v.path)-1
		if isLast && !v.append {
			// set value if map
			if cm, ok := cur.(map[string]any); ok {
				cm[part] = v.value
				return nil
			}
			// set value if slice
			if cs, ok := cur.([]any); ok {
				i, err := strconv.Atoi(part)
				if err != nil {
					return err
				}
				cs[i] = v.value
				return nil
			}
			fmt.Println(v.path, len(v.path))
			panic("unexpected value")
		}
		nextIsMap := true
		nextIndex := -1
		var err error
		if !isLast {
			nextIndex, err = strconv.Atoi(v.path[idx+1])
			if err == nil {
				nextIsMap = false
			}
		}
		if isLast && v.append {
			nextIsMap = false
		}
		switch c := cur.(type) {
		case map[string]any:
			// map key exists, set to current
			if _, ok := c[part]; ok {
				cur = c[part]
				cs, ok := cur.([]any)
				// if last part and is slice, do append
				if ok && isLast && v.append {
					c[part] = append(cs, v.value)
				}
				// if slice, and next part index is larger
				if ok && nextIndex >= len(cs) {
					needed := (nextIndex + 1) - len(cs)
					for i := 0; i < needed; i++ {
						c[part] = append(c[part].([]any), nil)
					}
					cur = c[part]
				}
				continue
			}
			// map key does not exist
			if nextIsMap {
				c[part] = make(map[string]any)
			} else {
				if v.append {
					c[part] = []any{v.value}
				} else {
					if nextIndex > -1 {
						c[part] = make([]any, nextIndex+1, nextIndex+1)
					} else {
						c[part] = []any{}
					}
				}
			}
			// set to current
			cur = c[part]
		case []any:
			i, err := strconv.Atoi(part)
			if err != nil {
				panic("non-numeric index for slice")
			}
			// slice index is valid/"exists", set to current
			if i < len(c) {
				cur = c[i]
				cs, ok := cur.([]any)
				// if last part and is slice, do append
				if ok && isLast && v.append {
					c[i] = append(cs, v.value)
				}
				// if slice, and next part index is larger
				if ok && nextIndex >= len(cs) {
					needed := (nextIndex + 1) - len(cs)
					for i := 0; i < needed; i++ {
						c[i] = append(cs, nil)
					}
				}
				// if nil, it doesn't "exist"
				if c[i] != nil {
					continue
				}
			}
			// slice index is "empty"
			if nextIsMap {
				c[i] = make(map[string]any)
			} else {
				if v.append {
					c[i] = []any{v.value}
				} else {
					if nextIndex > -1 {
						c[i] = make([]any, nextIndex+1, nextIndex+1)
					} else {
						c[i] = []any{}
					}
				}
			}
			// set to current
			cur = c[i]
		default:
			panic("unexpected type")
		}
	}
	return nil
}
