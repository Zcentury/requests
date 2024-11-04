package main

import (
	"github.com/Zcentury/requests"
	"github.com/Zcentury/requests/params"
)

func main() {
	url := "https://example.com"
	// 这两种写法都是发送json数据
	body := params.BodyJsonString(`{"name":"abc","age":18}`)
	body := params.BodyMap2Json{
		"name": "abc",
		"age":  "18",
	}

	// 这两种写法都是发送普通form-data数据
	body := params.BodyString(`name=abc&age=18`)
	body := params.BodyMap2String{
		"name": "abc",
		"age":  "18",
	}

	requests.Post(url, body)

}
