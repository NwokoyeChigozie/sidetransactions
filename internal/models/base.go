package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Tabler interface {
	TableName() string
}

type jsonmap map[string]interface{}
type roleCapabilities map[string]interface{}

// Value Marshal
func (a jsonmap) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Value Marshal
func (a roleCapabilities) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *jsonmap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		b, ok := value.(string)
		if !ok {
			return fmt.Errorf("type assertion to []byte or string failed: %v", value)
		}

		// Handle arrays
		var arr []interface{}
		if err := json.Unmarshal([]byte(b), &arr); err == nil {
			*a = make(map[string]interface{})
			(*a)["array"] = arr
			return nil
		}
		err := json.Unmarshal([]byte(b), &a)
		return err
	}

	// Handle arrays
	var arr []interface{}
	if err := json.Unmarshal(b, &arr); err == nil {
		*a = make(map[string]interface{})
		(*a)["array"] = arr
		return nil
	}

	return json.Unmarshal(b, &a)
}

func (a *roleCapabilities) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		b, ok := value.(string)
		if !ok {
			return fmt.Errorf("type assertion to []byte or string failed: %v", value)
		}

		// Handle arrays
		arr, err := handleArrayForRoleCapabilities([]byte(b))
		if err == nil {
			*a = arr
			return nil
		}

		err = json.Unmarshal([]byte(b), &a)
		return err
	}

	// Handle arrays
	arr, err := handleArrayForRoleCapabilities(b)
	if err == nil {
		*a = arr
		return nil
	}

	return json.Unmarshal(b, &a)
}

func handleArrayForRoleCapabilities(b []byte) (map[string]interface{}, error) {
	var arr []interface{}
	var a = map[string]interface{}{}

	err := json.Unmarshal(b, &arr)
	if err != nil {
		return a, err
	}

	for _, value := range arr {
		a[fmt.Sprintf("%v", value)] = true
	}
	return a, nil
}

func addQuery(oldQuery, newQuery, joinValue string) string {
	if oldQuery == "" {
		return " " + newQuery + " "
	}

	return oldQuery + " " + joinValue + " " + newQuery

}
