package utils

import (
	"strconv"
)

func QueryConvertInt(strQuery string, defaultVal int) (int, error) {
	newValue, err := strconv.ParseInt(strQuery, 10, 32)
	if err != nil {
		return defaultVal, err
	}
	return int(newValue), nil
}
