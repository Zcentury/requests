package main

import (
	"fmt"
	"github.com/Zcentury/requests"
	"github.com/Zcentury/requests/params"
	"log"
)

func main() {
	header := params.Headers{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
	}
	session := requests.Session()

	resp, err := session.Get("https://www.baidu.com/s?ie=UTF-8&wd=baidu", header)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Response.Cookies())

	resp1, err := session.Get("https://www.baidu.com/s?ie=UTF-8&wd=baidu", header)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(resp1.Response.Status)
	fmt.Println(resp1.Response.Cookies())

}
