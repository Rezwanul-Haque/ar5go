package methodsutil

import (
	"ar5go/infra/errors"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

func IsInvalid(value string) bool {
	if value == "" {
		return true
	}
	return false
}

// Contains search a needle(value) in the haystack(array)
func Contains(needle interface{}, haystack interface{}) bool {
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

func IsEmpty(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func MapToStruct(input map[string]interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

func StructToStruct(input interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

func ParseJwtToken(token, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidJwtSigningMethod
		}
		return []byte(secret), nil
	})
}

func StringToIntArray(stringArray []string) []int {
	var res []int

	for _, v := range stringArray {
		if i, err := strconv.Atoi(v); err == nil {
			res = append(res, i)
		}
	}

	return res
}
