package Common

import (
	"reflect"
	"strconv"
)

// Get List Tag of pointer of an struct
// Example: ListTag(&structVariable)
func ListTag(f interface{}) []string {
	val := reflect.ValueOf(f).Elem()

	listTag := []string{}

	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag
		listTag = append(listTag, string(tag))
	}

	return listTag
}

// Get List Tag of pointer of an struct
// Example: ListTag(&structVariable, "csv")
func ListTagByKey(f interface{}, tagKey string) []string {
	val := reflect.ValueOf(f).Elem()

	listTag := []string{}

	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get(tagKey)
		if len(tag) > 0 {
			listTag = append(listTag, tag)
		}
	}

	return listTag
}

// Parse struct value of pointer of an interface
// Example: ListTag(&structVariable)
func ListStringValue(f interface{}) []string {
	val := reflect.ValueOf(f).Elem()

	listTag := []string{}

	for i := 0; i < val.NumField(); i++ {
		v := ""
		switch val.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = strconv.FormatInt(val.Field(i).Int(), 10)
		default:
			v = val.Field(i).String()
		}
		listTag = append(listTag, v)
	}

	return listTag
}

// Parse struct value of pointer of an interface
// Example: ListTag(&structVariable, "csv")
func ListStringValueByTag(f interface{}, tagKey string) []string {
	val := reflect.ValueOf(f).Elem()

	listTag := []string{}

	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get(tagKey)
		if len(tag) > 0 {
			v := ""
			switch val.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v = strconv.FormatInt(val.Field(i).Int(), 10)
			default:
				v = val.Field(i).String()
			}
			listTag = append(listTag, v)
		}
	}

	return listTag
}
