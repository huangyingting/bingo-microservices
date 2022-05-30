package conf

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
)

func BsResolver(input map[string]interface{}) error {
	mapper := func(name string) string {
		args := strings.SplitN(strings.TrimSpace(name), ":", 2) //nolint:gomnd
		if v, has := readValue(input, args[0]); has {
			s, _ := v.String()
			return s
		} else if len(args) > 1 { // default value
			return args[1]
		}
		return ""
	}

	var resolve func(map[string]interface{}) error
	resolve = func(sub map[string]interface{}) error {
		for k, v := range sub {
			switch vt := v.(type) {
			case string:
				vs := expand(vt, mapper)
				if vst := strings.Trim(vs, "'"); len(vst) == len(vs)-1 {
					sub[k] = vst
				} else if vs == "true" || vs == "false" {
					vb, _ := strconv.ParseBool(vs)
					sub[k] = vb
				} else if vi, err := strconv.ParseInt(vs, 0, 32); err == nil {
					sub[k] = vi
				} else if vf, err := strconv.ParseFloat(vs, 32); err == nil {
					sub[k] = vf
				} else {
					sub[k] = vs
				}

			case map[string]interface{}:
				if err := resolve(vt); err != nil {
					return err
				}
			case []interface{}:
				for i, iface := range vt {
					switch it := iface.(type) {
					case string:
						vt[i] = expand(it, mapper)
					case map[string]interface{}:
						if err := resolve(it); err != nil {
							return err
						}
					}
				}
				sub[k] = vt
			}
		}
		return nil
	}
	return resolve(input)
}

// =============================================
// Copy from kratos and make no change

func expand(s string, mapping func(string) string) string {
	r := regexp.MustCompile(`\${(.*?)}`)
	re := r.FindAllStringSubmatch(s, -1)
	for _, i := range re {
		if len(i) == 2 { //nolint:gomnd
			s = strings.ReplaceAll(s, i[0], mapping(i[1]))
		}
	}
	return s
}

type atomicValue struct {
	atomic.Value
}

type ValueLite interface {
	String() (string, error)
	Store(interface{})
	Load() interface{}
}

func (v *atomicValue) String() (string, error) {
	switch val := v.Load().(type) {
	case string:
		return val, nil
	case bool, int, int32, int64, float64:
		return fmt.Sprint(val), nil
	case []byte:
		return string(val), nil
	default:
		if s, ok := val.(fmt.Stringer); ok {
			return s.String(), nil
		}
	}
	return "", fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.Load()))
}

// readValue read Value in given map[string]interface{}
// by the given path, will return false if not found.
func readValue(values map[string]interface{}, path string) (ValueLite, bool) {
	var (
		next = values
		keys = strings.Split(path, ".")
		last = len(keys) - 1
	)
	for idx, key := range keys {
		value, ok := next[key]
		if !ok {
			return nil, false
		}
		if idx == last {
			av := &atomicValue{}
			av.Store(value)
			return av, true
		}
		switch vm := value.(type) {
		case map[string]interface{}:
			next = vm
		default:
			return nil, false
		}
	}
	return nil, false
}
