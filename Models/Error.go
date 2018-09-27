package Models

import "github.com/ruslanfedoseenko/dhtcrawler/Errors"

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewError(code int) Error {
	return Error{
		Code:    code,
		Message: Errors.ErrorText(code),
	}
}

func NewErrorAddText(code int, e error) Error {
	return Error{
		Code:    code,
		Message: Errors.ErrorText(code) + " " + e.Error(),
	}
}
