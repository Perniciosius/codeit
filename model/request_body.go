package model

type CompileRequestBody struct {
	Code      string `json:"code"`
	Arguments string `json:"arguments"`
}
