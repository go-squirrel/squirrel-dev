package res

import "squirrel-dev/internal/pkg/response"

const (
	ErrInvalidScriptContent = 80001
	ErrScriptNotFound       = 80002
	ErrDuplicateScript      = 80003
	ErrUnsupportedScriptType = 80004
)

func RegisterCode() {
	response.Register(ErrInvalidScriptContent, "invalid script content")
	response.Register(ErrScriptNotFound, "script not found")
	response.Register(ErrDuplicateScript, "script already exists")
	response.Register(ErrUnsupportedScriptType, "unsupported script type")
}
