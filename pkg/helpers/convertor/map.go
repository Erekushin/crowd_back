package convertor

import (
	"encoding/json"
	"errors"
)

// MapToMap :
func MapToMap(param map[string]interface{}, key, errorString string) (result map[string]interface{}, err error) {
	result, ok := param[key].(map[string]interface{})
	if !ok {
		return result, errors.New(errorString)
	}
	return result, err
}

// MapToStruct : map to struct
func MapToStruct(m interface{}, v interface{}) (err error) {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return err
	}
	json.Unmarshal(jsonString, v)
	return err
}

// MapToKeyValueSlice :
func MapToKeyValueSlice(m, e map[string]interface{}) (keys, values []string) {
	for key, value := range m {
		if _, ok := e[key]; ok {
			continue
		}
		vs, err := InterfaceToString(value, "err")
		if err == nil {
			keys = append(keys, key)
			values = append(values, vs)
		}
	}
	return keys, values
}
