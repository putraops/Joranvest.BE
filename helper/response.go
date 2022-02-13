package helper

import (
	"strings"
)

type JSONResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type EmptyObj struct{}

func BuildResponse(status bool, message string, data interface{}) JSONResponse {
	res := JSONResponse{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, err string, data interface{}) JSONResponse {
	splittedErros := strings.Split(err, "\n")
	res := JSONResponse{
		Status:  false,
		Message: message,
		Errors:  splittedErros,
		Data:    data,
	}
	return res
}

func ServerResponse(status bool, message string, err string, data interface{}) Response {
	splittedErros := strings.Split(err, "\n")
	res := Response{
		Status:  status,
		Message: message,
		Errors:  splittedErros,
		Data:    data,
	}
	return res
}
