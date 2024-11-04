package requests

import (
	"io"
	"net/http"
)

type Response struct {
	Response   *http.Response // 原始 http.Response
	Text       string         // 字符串Body
	StatusCode int            // 状态码
	Body       io.ReadCloser  // Copy过来的Body流，流操作用这个，因为Response的Body已经关闭了
	Header     http.Header    // 响应头
}
