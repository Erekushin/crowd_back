package convertor

import (
	"encoding/json"
	"strconv"
	"time"
)

// StringToInt
func StringToInt(num string) int {
	res, err := strconv.Atoi(num)
	if err != nil {
		res = 0
	}
	return res
}

// StringToMap : convert string to map string interface
func StringToMap(jsonString string) (result map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(jsonString), &result)
	return result, err
}

// StringToMapArr : convert string to map string interface array
func StringToMapArr(jsonString string) (result []map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(jsonString), &result)
	return result, err
}

// DateStringToTime
func DateStringToTime(t string) time.Time {
	c, err := time.Parse("2006-01-02", t)
	if err != nil {
		c, _ = time.Parse("2006-01-02", "0000-00-00")
	}

	return c
}
