package requests

import (
	"io"
	"net/http"
)

type Response struct {
	Response   *http.Response
	Text       string
	StatusCode int
	Body       io.ReadCloser
	Header     http.Header
}
