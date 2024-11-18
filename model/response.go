package model

type ResponseInfo struct {
	Code int32       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
