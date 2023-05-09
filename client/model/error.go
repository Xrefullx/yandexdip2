package model

import "errors"

var (
	ErrorItemNotFound  = errors.New("item not found")
	ErrorParamNotValid = errors.New("incoming parameter not valid")
)
