package core

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Unmarshal dynamically unmarshals JSON string into the provided entity type
func Unmarshal(args string, entityType interface{}) (interface{}, error) {
	entityValue := reflect.New(reflect.TypeOf(entityType)).Interface()
	err := json.Unmarshal([]byte(args), entityValue)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json string: %v", err)
	}
	return entityValue, nil
}
