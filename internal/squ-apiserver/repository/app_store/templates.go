package app_store

import "embed"

//go:embed templates/*.yml
var TemplateFS embed.FS
