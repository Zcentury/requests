package requests

import (
	"github.com/Zcentury/requests/method"
)

func Request(method method.Method, url string, args ...interface{}) (*Response, error) {
	session := Session()
	return session.Request(method, url, args...)
}

func Get(url string, args ...interface{}) (*Response, error) {
	return Request(method.GET, url, args...)
}

func Post(url string, args ...interface{}) (*Response, error) {
	return Request(method.POST, url, args...)
}

func Put(url string, args ...interface{}) (*Response, error) {
	return Request(method.PUT, url, args...)
}

func Delete(url string, args ...interface{}) (*Response, error) {
	return Request(method.DELETE, url, args...)
}
