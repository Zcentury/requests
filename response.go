package requests

import "net/http"

type Response struct {
	Response *http.Response
	Text     string
}
