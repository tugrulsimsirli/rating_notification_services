package utils

// Simple AutoMapper-like utility to map struct fields
import (
	"reflect"
)

// Map maps data from one struct to another
func Map(src interface{}, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		field := srcVal.Type().Field(i)
		dstField := dstVal.FieldByName(field.Name)
		if dstField.IsValid() && dstField.CanSet() {
			dstField.Set(srcVal.Field(i))
		}
	}
}
