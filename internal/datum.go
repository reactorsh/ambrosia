package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type datum map[string]interface{}

func (d datum) String(keys []string, fields bool) string {
	var str strings.Builder
	for _, key := range keys {
		if value, ok := d[key]; ok {
			if fields {
				str.WriteString(fmt.Sprintf("%s: %v\n", key, value))
			} else {
				str.WriteString(fmt.Sprintf("%v\n", value))
			}
		}
	}
	return str.String()
}

func (d datum) JSON(keys []string) ([]byte, error) {
	selected := make(datum)

	for _, key := range keys {
		if val, ok := d[key]; ok {
			selected[key] = val
		}
	}

	bytes, err := json.Marshal(selected)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func isEqual(a datum, b datum) bool {
	delete(a, "ambrosia")
	delete(b, "ambrosia")
	return reflect.DeepEqual(a, b)
}

func contains(s []datum, e datum) bool {
	for _, a := range s {
		if isEqual(a, e) {
			return true
		}
	}
	return false
}

func datumSub(a []datum, b []datum) []datum {
	var ret []datum
	for _, item := range a {
		if !contains(b, item) {
			ret = append(ret, item)
		}
	}
	return ret
}
