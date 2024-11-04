package requests

import "errors"

var (
	ErrNoUrl  = errors.New("请传入URL")
	ErrNoBody = errors.New("请传入Body")
)
