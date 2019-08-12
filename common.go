package httpreq

import (
	"encoding/base64"
	"encoding/json"
)

func MarshalString(i interface{}) string {
	data, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(data)
}

func Marshal(i interface{}) []byte {
	data, _ := json.Marshal(i) //nolint errcheck
	return data
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
