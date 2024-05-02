package util

import "reflect"

// Contains checks if a value is present in a slice using reflection.
func Contains(slice interface{}, value interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		panic("Contains expects a slice as the first argument")
	}

	for i := 0; i < sliceValue.Len(); i++ {
		element := sliceValue.Index(i).Interface()
		if reflect.DeepEqual(element, value) {
			return true
		}
	}
	return false
}
