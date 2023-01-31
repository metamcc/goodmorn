package tooltip

import (
	"time"
	"encoding/binary"
	"math"
)


func FormatDateUTC2LTZ(d string, formatSRC string, formatTGT string, timezoneOffsetSeconds int) string {

	dateUTC, _ := time.Parse(formatSRC, d)

	offset := time.Duration(timezoneOffsetSeconds)*time.Second
	dateLTZ := dateUTC.Add(offset)

	return string(dateLTZ.Format(formatTGT))
}

func FormatDate(d string, formatSRC string, formatTGT string) string {
	date, _ := time.Parse(formatSRC, d)

	return string(date.Format(formatTGT))
}

func ParseDate(d string, dateFormat string) time.Time {
	date, _ := time.Parse(dateFormat, d)
	return date
}

func CalculateDaysBetween(endDate string, startDate string, dateFormat string) int {

	date1, _ := time.Parse(dateFormat, endDate)
	date2, _ := time.Parse(dateFormat, startDate)
	difference := date1.Sub(date2)
	daysBetween := int(difference.Hours() / 24)

	return daysBetween
}

func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
