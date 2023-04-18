package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONRaw json.RawMessage

func (json JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(json)
	return driver.Value(byteArr), nil
}

func (j *JSONRaw) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan хийж буй зүйл нь bytes төрөл биш байна"))
	}
	err := json.Unmarshal(asBytes, &j)
	if err != nil {
		return error(errors.New("Scan unmarshal хийхэд алдаа гарлаа"))
	}
	return nil
}

func (j *JSONRaw) MarshalJSON() ([]byte, error) {
	return *j, nil
}

func (j *JSONRaw) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: Unmarshal хийх элемент нь nil байна")
	}
	*j = append((*j)[0:0], data...)
	return nil
}
