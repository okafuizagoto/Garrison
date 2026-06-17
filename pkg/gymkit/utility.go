package gymkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
)

// GetEnv returns env value or fallback if not set.
func GetEnv(key string, fallback string) string {
	env := os.Getenv(key)
	if len(env) == 0 {
		return fallback
	}
	return env
}

// IsJSON checks whether a string is valid JSON object.
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// IsJSONString checks whether a string is a valid JSON string value.
func IsJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil
}

// IsSuccessCode returns true for 2xx HTTP status codes.
func IsSuccessCode(c int) bool {
	return c > 199 && c < 300
}

// IsErrorCode returns true for 4xx+ HTTP status codes.
func IsErrorCode(c int) bool {
	return c > 399
}

// InArray checks whether needle exists in haystack slice/array.
func InArray(needle interface{}, haystack interface{}) (bool, int, error) {
	haystackValue := reflect.ValueOf(haystack)
	if haystackValue.Kind() != reflect.Array && haystackValue.Kind() != reflect.Slice {
		return false, -1, errors.New("parameter 2 is not an array or slice")
	}
	for i := 0; i < haystackValue.Len(); i++ {
		if reflect.DeepEqual(haystackValue.Index(i).Interface(), needle) {
			return true, i, nil
		}
	}
	return false, -1, nil
}

// IsEmptyString returns true if str contains only whitespace.
func IsEmptyString(str string) bool {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(str, "") == ""
}

// ToString converts common types to string.
func ToString(input interface{}) string {
	switch v := input.(type) {
	case string:
		return v
	case map[string]interface{}:
		b, _ := json.Marshal(v)
		return string(b)
	}
	return fmt.Sprintf("%v", input)
}
