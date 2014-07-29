package internal

import (
	"reflect"
	"strconv"
)

func InspectKind(str string) reflect.Kind {
	switch {
	case str == "true" || str == "false":
		return reflect.Bool
	case isInt64(str):
		return reflect.Int64
	case isFloat64(str):
		return reflect.Float64
	default:
		return reflect.String
	}
}

func isInt64(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func isFloat64(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
