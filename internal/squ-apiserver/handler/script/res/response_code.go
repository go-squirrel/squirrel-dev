package res

import "squirrel-dev/internal/pkg/response"

const (
	// 80000-80019: 基础操作
	ErrScriptNotFound        = 80001
	ErrDuplicateScript       = 80002
	ErrInvalidScriptContent  = 80003
	ErrUnsupportedScriptType = 80004
	ErrScriptNotDeployed     = 80005

	// 80020-80039: 执行相关
	ErrScriptExecutionFailed = 80021
	ErrScriptTimeout         = 80022
	ErrServerNotFound        = 80023
)

func RegisterCode() {
	// 80000-80019: 基础操作
	response.Register(ErrScriptNotFound, "script not found")
	response.Register(ErrDuplicateScript, "script already exists")
	response.Register(ErrInvalidScriptContent, "invalid script content")
	response.Register(ErrUnsupportedScriptType, "unsupported script type")
	response.Register(ErrScriptNotDeployed, "script not deployed")

	// 80020-80039: 执行相关
	response.Register(ErrScriptExecutionFailed, "script execution failed")
	response.Register(ErrScriptTimeout, "script execution timeout")
	response.Register(ErrServerNotFound, "server not found")
}
