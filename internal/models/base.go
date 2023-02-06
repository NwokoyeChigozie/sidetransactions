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

// Value Marshal
func (a jsonmap) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *jsonmap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		b, ok := value.(string)
		if !ok {
			return fmt.Errorf("type assertion to []byte or string failed: %v", value)
		}
		return json.Unmarshal([]byte(b), &a)
	}
	return json.Unmarshal(b, &a)
}
