package static

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DefaultConfig returns a generic default configuration for SPA middleware.
func DefaultConfig() Config {
	return Config{
		IndexFilePath: "index.html",
		APIPrefix:     "/api",
		StaticPaths:   []string{"/assets", "/_nuxt", "/favicon.ico"},
	}
}

// Default returns the SPA middleware with default configuration.
func Default(staticFS embed.FS, distDir string) gin.HandlerFunc {
	config := DefaultConfig()
	return New(staticFS, distDir, config)
}

// New returns the SPA middleware with user-defined custom configuration.
func New(staticFS embed.FS, distDir string, config Config) gin.HandlerFunc {
	// Set defaults
	if config.IndexFilePath == "" {
		config.IndexFilePath = DefaultConfig().IndexFilePath
	}
	if config.APIPrefix == "" {
		config.APIPrefix = DefaultConfig().APIPrefix
	}
	if len(config.StaticPaths) == 0 {
		config.StaticPaths = DefaultConfig().StaticPaths
	}

	// Build full path to index.html
	indexPath := distDir + "/" + config.IndexFilePath

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Check if this is an API request
		if strings.HasPrefix(path, config.APIPrefix) {
			c.Next()
			return
		}

		// Check if this is a static resource request
		for _, staticPath := range config.StaticPaths {
			if strings.HasPrefix(path, staticPath) {
				c.Next()
				return
			}
		}

		// All other requests should return index.html
		c.Header("Content-Type", "text/html; charset=utf-8")

		// Read index.html from embedded file system
		indexContent, err := fs.ReadFile(&staticFS, indexPath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load index.html")
			c.Abort()
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", indexContent)
		c.Abort()
	}
}
