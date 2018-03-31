package goconcept

import (
	"net/url"
	"strconv"
)

func Util__queryToInt(vars url.Values, var_name string, min int, max int, min_inf bool, max_inf bool, default_value int) int {
	val_str, ok := vars[var_name]
	if !ok {
		return default_value
	}
	if len(val_str) < 1 {
		return default_value
	}
	val_int, err := strcon.Atoi(val_str[0])
	if !min_inf && val_int < min {
		return default_value
	}
	if !max_inf && val_int > max {
		return default_value
	}
	return val_int
}
