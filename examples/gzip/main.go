package main

import (
	"fmt"
	"github.com/Zcentury/requests"
	"github.com/Zcentury/requests/params"
)

func main() {
	headers := params.Headers{
		"Accept-Encoding": "gzip, deflate",
	}
	response, err := requests.Get("https://example.com", headers)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Text)

}
