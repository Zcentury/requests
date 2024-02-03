package requests

import (
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
