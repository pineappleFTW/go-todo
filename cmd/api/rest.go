package main

var (
	Success = "success"
	Fail    = "fail"
	Error   = "error"
)

type response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code"`
}
