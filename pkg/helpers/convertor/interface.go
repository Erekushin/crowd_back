package convertor

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

// InterfaceToBytes
func InterfaceToBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

// InterfaceGetBytes : interface{} to bytes
func InterfaceGetBytes(key interface{}) ([]byte, error) {
	b, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// InterfaceToInt64
func InterfaceToInt64(param interface{}, errorString string) (result int64, err error) {
	result, ok := param.(int64)
	if !ok {
		if fl, ok := param.(float64); ok {
			result = int64(fl)
		} else {
			return result, errors.New(errorString)
		}
	}
	return result, err
}

// InterfaceToString :
func InterfaceToString(param interface{}, errorString string) (result string, err error) {
	result, ok := param.(string)
	if !ok {
		return result, errors.New(errorString)
	}
	return result, err
}

// InterfaceToMap :

func InterfaceToMap(param interface{}) (result map[string]interface{}, err error) {
	inrec, err := json.Marshal(param)
	json.Unmarshal(inrec, &result)
	return result, err
}
