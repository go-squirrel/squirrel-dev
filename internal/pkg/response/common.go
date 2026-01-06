package response

import "sync"

var (
	codeMsgMap map[int]string
	once       sync.Once
	mu         sync.RWMutex // 用于并发安全（可选，但推荐）
)

// 错误码
const (
	CodeSuccess = 0

	ErrCodeParameter  = 41001
	ErrUserOrPassword = 41002

	ErrSQL           = 50000
	ErrSQLNotFound   = 50001
	ErrSQLNotUnique  = 50002
	ErrDuplicatedKey = 50003
)

func Init() {
	once.Do(func() {
		codeMsgMap = make(map[int]string)
		baseRes(codeMsgMap)
	})
}

func baseRes(msg map[int]string) map[int]string {
	msg[CodeSuccess] = "success"

	msg[ErrCodeParameter] = "parameter error"
	msg[ErrUserOrPassword] = "user or password error"

	msg[ErrSQL] = "sql error"
	msg[ErrSQLNotFound] = "sql not found"
	msg[ErrSQLNotUnique] = "sql not unique"
	msg[ErrDuplicatedKey] = "duplicated key"
	return msg
}

func getMessage(code int) (message string) {
	message, ok := codeMsgMap[code]
	if !ok {
		message = "unknown error"
	}
	return message
}
