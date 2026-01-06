package response

var codeMsgMap map[int]string

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
	codeMsgMap = make(map[int]string, 1024)
	codeMsgMap = baseRes(codeMsgMap)
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
