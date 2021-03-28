package model

type CompileRequestBody struct {
	Code             string `json:"code" form:"code"`
	Input            string `json:"input" form:"input"`
	CompileArguments string `json:"compileArguments" form:"compileArguments"`
	RuntimeArguments string `json:"runtimeArguments" form:"runtimeArguments"`
}
