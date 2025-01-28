package main

import "encoding/json"

type RpcCall struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	Id      uint            `json:"id"`
}

type RpcResponse struct {
	Id     uint           `json:"id"`
	Reuslt any            `json:"result"`
	Error  *ResponseError `json:"error"`
}

type ResponseError struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
