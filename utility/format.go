package utility

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func FormatDate(date, currentISOFormat, newISOFormat string) (string, error) {
	t, err := time.Parse(currentISOFormat, date)
	if err != nil {
		return date, err
	}
	return t.Format(newISOFormat), nil
}

func UnFormatDueDate(date string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		t, err = time.Parse("2006-01-02 15::05", date)
		if err != nil {
			return t, err
		}
	}
	return t, nil
}
func FormatDateSpecialCase(t time.Time) string {
	day := t.Day()
	var suffix string
	switch {
	case day%10 == 1 && day != 11:
		suffix = "st"
	case day%10 == 2 && day != 12:
		suffix = "nd"
	case day%10 == 3 && day != 13:
		suffix = "rd"
	default:
		suffix = "th"
	}
	return t.Format(fmt.Sprintf("2%s January, 2006", suffix))
}

func GetUnixTime(date, currentISOFormat, newISOFormat string) (int, error) {
	t, err := time.Parse(currentISOFormat, date)
	if err != nil {
		return 0, err
	}
	return int(t.Unix()), nil
}
func GetUnixString(date, currentISOFormat, newISOFormat string) (string, error) {
	t, err := time.Parse(currentISOFormat, date)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(t.Unix())), nil
}

func GetStartAndEnd(interval string) (time.Time, time.Time) {
	var (
		start time.Time
		end   time.Time
	)
	switch interval {
	case "day":
		startOfDay := time.Now().Truncate(24 * time.Hour)
		start = startOfDay
		end = startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)
	case "monday_to_thursday":
		now := time.Now()
		startOfWeek := now.AddDate(0, 0, -int(now.Weekday()-1))
		start = startOfWeek
		end = startOfWeek.AddDate(0, 0, 4)
	case "week":
		startOfWeek := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
		start = startOfWeek
		end = startOfWeek.Add(24*7*time.Hour - time.Nanosecond)
	case "month":
		startOfMonth := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -int(time.Now().Day()-1))
		start = startOfMonth
		end = startOfMonth.AddDate(0, 1, -1).Add(24*time.Hour - time.Nanosecond)
	default:
		start = time.Now()
		end = time.Now()
	}
	return start, end
}

func ConvertStringInterfaceToStringFloat(originalMap map[string]interface{}) map[string]float64 {
	convertedMap := make(map[string]float64)
	for key, value := range originalMap {
		if val, ok := value.(float64); ok {
			convertedMap[key] = val
		} else if val, ok := value.(string); ok {
			if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
				convertedMap[key] = floatVal
			}
		}
	}
	return convertedMap
}

func RemoveKey(p interface{}, key string) {
	val := reflect.ValueOf(p).Elem()
	val.FieldByName(key).Set(reflect.Zero(val.FieldByName(key).Type()))
}

func CopyStruct(src, dst interface{}) {
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)

		if dstField.IsValid() {
			dstField.Set(srcField)
		}
	}
}
