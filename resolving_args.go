package requests

import (
	"encoding/json"
	"github.com/Zcentury/requests/params"
	urlutil "net/url"
	"reflect"
)

// 解析参数
func resolvingArgs(args ...interface{}) (map[string]interface{}, string, error) {
	result := make(map[string]interface{})
	var contentType string

	for _, arg := range args {

		switch v := reflect.TypeOf(arg); v.Kind() {

		case reflect.String:
			switch v.String() {
			case "params.UrlParams":
				result["UrlParams"] = string(arg.(params.UrlParams))
			case "params.BodyJsonString":
				result["Body"] = string(arg.(params.BodyJsonString))
				contentType = "application/json"
			case "params.BodyString":
				result["Body"] = string(arg.(params.BodyString))
				contentType = "application/x-www-form-urlencoded"

			}

		case reflect.Map:

			switch v.String() {
			case "params.Headers":
				result["Headers"] = arg.(params.Headers)

			case "params.UrlMap2Params":
				values := urlutil.Values{}
				for key, value := range arg.(params.UrlMap2Params) {
					values.Add(key, value)
				}
				result["UrlParams"] = values.Encode()

			case "params.BodyMap2Json":
				if jsonData, err := json.Marshal(arg.(params.BodyMap2Json)); err == nil && result["Body"] != "" {
					result["Body"] = string(jsonData)
				}
				contentType = "application/json"

			case "params.BodyMap2String":
				values := urlutil.Values{}
				for key, value := range arg.(params.BodyMap2String) {
					values.Add(key, value)
				}
				result["Body"] = values.Encode()
				contentType = "application/x-www-form-urlencoded"

			case "params.Proxy":
				result["Proxy"] = arg.(params.Proxy)

			}

		}

	}
	return result, contentType, nil
}
