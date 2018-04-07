package numbers

import (
	"encoding/binary"
	"math"
	"strconv"
)

// ByteToFloat64 conversion
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

// Float64ToByte conversion
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

// StringToFloat64 conversion
func StringToFloat64(value string) (float64, error) {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Float64ToString conversion
func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 6, 64)
}

// ToInt64 any byte array
func ToInt64(b []byte) int64 {
	number, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0
	}
	return int64(number)
}

// Int64ToString cast
func Int64ToString(number int64) string {
	return strconv.FormatInt(number, 10)
}

// Int64ToByte cast
func Int64ToByte(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

// ByteToInt64 cast
func ByteToInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}
