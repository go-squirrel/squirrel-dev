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

	ErrCodeParameter     = 41001
	ErrUserOrPassword    = 41002
	ErrMissingAuthHeader = 41003
	ErrInvalidAuthFormat = 41004
	ErrTokenInvalid      = 41005

	ErrTLSRequired       = 41006
	ErrMissingClientCert = 41007
	ErrCertCNNotAllowed  = 41008
	ErrCertVerifyFailed  = 41009

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
	msg[ErrMissingAuthHeader] = "missing authorization header"
	msg[ErrInvalidAuthFormat] = "invalid authorization header format"
	msg[ErrTokenInvalid] = "invalid or expired token"
	msg[ErrTLSRequired] = "tls connection required"
	msg[ErrMissingClientCert] = "missing client certificate"
	msg[ErrCertCNNotAllowed] = "client certificate common name not allowed"
	msg[ErrCertVerifyFailed] = "client certificate verification failed"

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
