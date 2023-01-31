package tooltip

import (
	"encoding/binary"
	"math"
	"regexp"
	"strconv"
	"time"
)

func IsNumeric(s []string) bool {
	for i := 0; i < len(s); i++ {
		for _, v := range s[i] {
			if v < '0' || v > '9' {
				return false
			}
		}
	}
	return true
}

func CheckIfAllEqual(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[0] {
			return false
		}
	}
	return true
}

func GetNumbOfIterates(symbToFind string, strToCheck string) int {

	argSep := regexp.MustCompile(symbToFind)

	matches := argSep.FindAllStringIndex(strToCheck, -1)

	iNumb := len(matches)
	return iNumb
}

/* CheckExchangePeriod function to calculate difference in days between 2 dates
arg[0] : today_date
arg[1] : check_date
*/
func CheckExchangePeriod(date1 string, date2 string) int {

	format := "2006-01-02 15:04:05"
	modelDate, _ := time.Parse(format, date1)
	checkDate, _ := time.Parse(format, date2)
	diff := modelDate.Sub(checkDate)
	daydiff := int(diff.Seconds())
	return daydiff
}

func FormatDateUTC2LTZ(d string, formatSRC string, formatTGT string, timezoneOffsetSeconds int) string {

	dateUTC, _ := time.Parse(formatSRC, d)

	offset := time.Duration(timezoneOffsetSeconds) * time.Second
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

func AddDaysToDate(date string, mints string) string {
	dateFormat := "2006-01-02 15:04:05"
	startDate, _ := time.Parse(dateFormat, date)

	minutsNumb, _ := strconv.Atoi(mints)

	endDate := startDate.Add(time.Minute * time.Duration(minutsNumb))
	EndDateToString := endDate.String()

	return EndDateToString[:19]
}

//formatDate function brngs date to "yyyy-mm-dd 00:00:00" format
func formatDate(d string) string {
	tempDate := d[:11]
	date := tempDate + "00:00:00"
	return date
}

func CalculateDaysBetween(startDate string, endDate string) int {

	dateFormat := "2006-01-02"
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
