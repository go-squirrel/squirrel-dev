package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func GenAgentUrl(schema, address string, port int, baseuri, uri string) string {
	host := address
	if port != 0 {
		host = fmt.Sprintf("%s:%d", address, port)
	}

	// 合并路径，确保只有一个斜杠分隔
	path := strings.TrimSuffix(baseuri, "/") + "/" + strings.TrimPrefix(uri, "/")
	path = strings.TrimPrefix(path, "/") // 确保 path 不以 / 开头，因为 url.Path 会自动加

	u := url.URL{
		Scheme: schema,
		Host:   host,
		Path:   path,
	}

	return u.String()
}
