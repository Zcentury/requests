package main

import (
	"fmt"
	"github.com/Zcentury/requests"
	"github.com/Zcentury/requests/params"
)

func main() {
	// 三选一，暂时不支持同时设置多个
	proxy := params.Proxy{
		"http": "http://127.0.0.1:8080",
		// "https":  "https://127.0.0.1:8080",
		// "socks5": "127.0.0.1:7890",
	}
	response, err := requests.Get("https://example.com", proxy)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Text)
}
