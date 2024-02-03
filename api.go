package requests

import (
	"github.com/Zcentury/requests/method"
)

func Request(method method.Method, url string, args ...interface{}) *Response {
	session := Session()
	return session.Request(method, url, args...)
}

func Get(url string, args ...interface{}) *Response {
	return Request(method.GET, url, args...)
}

func Post(url string, args ...interface{}) *Response {
	//pc, file, line, _ := runtime.Caller(1)
	//// 获取函数名
	//funcName := runtime.FuncForPC(pc).Name()
	//
	//fmt.Printf("Function: %s\n", funcName)
	//fmt.Printf("File: %s\n", file)
	//fmt.Printf("Line: %d\n", line)

	return Request(method.POST, url, args...)
}

//func Put(url string, args ...interface{}) *Response {
//	return Request(method.PUT, url, args...)
//}
