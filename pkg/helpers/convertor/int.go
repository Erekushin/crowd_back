package convertor

import "strconv"

func UintToString(num uint) string {
	return strconv.FormatUint(uint64(num), 10)
}
