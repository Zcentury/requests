package requests

import (
	"github.com/Zcentury/requests/method"
	"net/http"
	"net/http/cookiejar"
)

func Session() *Requests {
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}
	return NewRequests(client)
}

func (r *Requests) Get(url string, args ...interface{}) (*Response, error) {
	return r.Request(method.GET, url, args...)
}

func (r *Requests) Post(url string, args ...interface{}) (*Response, error) {
	return r.Request(method.POST, url, args...)
}

func (r *Requests) Put(url string, args ...interface{}) (*Response, error) {
	return r.Request(method.PUT, url, args...)
}

func (r *Requests) Delete(url string, args ...interface{}) (*Response, error) {
	return r.Request(method.DELETE, url, args...)
}
