package util

import "encoding/json"

func ToJsonStr(t interface{}) string {
	s, _ := json.Marshal(t)
	return string(s)
}
