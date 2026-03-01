package utils

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"
)

// GenAgentUrl constructs a full URL from components.
// - If address already contains a port (e.g., "host:port" or "[::1]:80"), the provided port is ignored.
// - If address has no port and port != 0, the port is appended.
// - Handles IPv4, IPv6 (with brackets), and domain names correctly.
// - Properly joins baseuri and uri using path.Join for robust path concatenation.
// - Supports query string in uri (e.g., "monitor/base?range=1h")
func GenAgentUrl(schema, address string, port int, baseuri, uri string) string {
	host := address

	// Check if address already includes a port
	if _, _, err := net.SplitHostPort(address); err != nil {
		// address does NOT contain a port
		if port != 0 {
			// Check if address is IPv6 without brackets
			if ip := net.ParseIP(address); ip != nil && ip.To4() == nil {
				// It's an IPv6 address → wrap in brackets
				host = fmt.Sprintf("[%s]:%d", address, port)
			} else {
				// Regular domain or IPv4
				host = fmt.Sprintf("%s:%d", address, port)
			}
		}
		// else: keep host = address (no port added)
	}
	// else: address already has port, use as-is

	// Separate path and query string from uri
	pathPart := uri
	queryPart := ""
	if idx := strings.Index(uri, "?"); idx != -1 {
		pathPart = uri[:idx]
		queryPart = uri[idx+1:]
	}

	// Normalize and join paths
	// Ensure both parts are treated as path segments
	fullPath := path.Join("/", baseuri, pathPart)
	// path.Join always returns absolute path (starts with /)

	u := url.URL{
		Scheme:   schema,
		Host:     host,
		Path:     fullPath,
		RawQuery: queryPart,
	}

	return u.String()
}
