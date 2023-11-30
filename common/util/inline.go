package util

import "reflect"

func MapKeys(m interface{}) interface{} {
	v := reflect.ValueOf(m)
	vKeys := v.MapKeys()
	mapKeyType := reflect.TypeOf(m).Key()
	retSliceType := reflect.SliceOf(mapKeyType)
	retSliceValue := reflect.MakeSlice(retSliceType, 0, 0)
	for _, v := range vKeys {
		retSliceValue = reflect.Append(retSliceValue, v)
	}
	return retSliceValue.Interface()
}

func MapValues(m interface{}) interface{} {
	mapValue := reflect.ValueOf(m)
	vKeys := mapValue.MapKeys()
	mapValType := reflect.TypeOf(m).Elem()
	retSliceType := reflect.SliceOf(mapValType)
	retSliceValue := reflect.MakeSlice(retSliceType, 0, 0)
	for _, v := range vKeys {
		retSliceValue = reflect.Append(retSliceValue, mapValue.MapIndex(v))
	}
	return retSliceValue.Interface()
}

func MapKeysString(m interface{}) []string {
	return MapKeys(m).([]string)
}

func MapValuesString(m interface{}) []string {
	return MapValues(m).([]string)
}
