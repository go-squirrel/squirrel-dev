package static

// Config represents all available options for the SPA middleware.
type Config struct {
	// IndexFilePath is the path to the index.html file within the dist directory
	// Default value is "index.html"
	IndexFilePath string

	// APIPrefix is the prefix for API routes that should not be handled by the SPA middleware
	// Default value is "/api"
	APIPrefix string

	// StaticPaths is a list of paths that should be served as static files
	// Default values are ["/assets", "/_nuxt", "/favicon.ico"]
	StaticPaths []string
}
