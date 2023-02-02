package utility

import (
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
