package response

// Response 回调时的固定内容
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error 回调错误信息
func Error(code int) Response {
	return Response{
		Code:    code,
		Message: getMessage(code),
	}
}

func ErrorUnknown(code int, data string) Response {
	return Response{
		Code:    code,
		Message: data,
	}

}

// Success 回调正确信息
func Success(data any) Response {
	return Response{
		Code:    CodeSuccess,
		Message: getMessage(CodeSuccess),
		Data:    data,
	}
}

func Register(code int, message string, force ...bool) {
	Init() // 确保已初始化

	doForce := len(force) > 0 && force[0]

	mu.Lock()
	defer mu.Unlock()

	if _, exists := codeMsgMap[code]; !exists || doForce {
		codeMsgMap[code] = message
	}
}
