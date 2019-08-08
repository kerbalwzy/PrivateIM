package models

import (
	"fmt"
	"log"
	"reflect"
)

// get a field value if a struct by the field name
func GetReflectValueByField(instance interface{}, field string) (interface{}, error) {
	fieldNames := GetReflectFieldNames(instance)
	if _, ok := fieldNames[field]; !ok {
		return nil, fmt.Errorf("getReflectValueByFeild fail :field <%s> dose not exist", field)
	}
	v := reflect.ValueOf(instance)
	return v.FieldByName(field).Interface(), nil
}

// get all field name of a struct
func GetReflectFieldNames(structName interface{}) map[string]int {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make(map[string]int)
	for i := 0; i < fieldNum; i++ {
		result[t.Field(i).Name] = i
	}
	return result
}
