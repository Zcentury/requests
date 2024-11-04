package main

import (
	"fmt"
	"github.com/Zcentury/requests"
)

func main() {
	response, err := requests.Get("https://example.com")
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Text)
}
