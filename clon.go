package clon

import (
	"encoding/json"
	"strconv"
	"strings"
)

func Parse(args []string) (any, error) {
	var pargs []arg
	for _, a := range args {
		pargs = append(pargs, parseArg(a))
	}
	if len(pargs[0].path) == 0 {
		var v []interface{}
		for _, a := range pargs {
			v = append(v, a.value)
		}
		return v, nil
	}
	v := make(map[string]interface{})
	for _, a := range pargs {
		v[strings.Join(a.path, ".")] = a.value
	}
	return v, nil
}

type arg struct {
	path    []string
	value   any
	literal bool
	append  bool
}

func parseArg(s string) (a arg) {
	kv := strings.Split(s, "=")
	path := kv[0]

	if strings.HasSuffix(path, ":") {
		a.literal = true
		path = path[:len(path)-1]
		a.value = parseLiteral(kv[1])
	} else {
		a.value = strings.Trim(kv[1], "'")
	}
	if strings.HasSuffix(path, "[]") {
		a.append = true
		path = path[:len(path)-2]
	}
	path = strings.ReplaceAll(path, "]", "")
	a.path = strings.Split(path, "[")
	return a
}

func parseLiteral(v string) any {
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
