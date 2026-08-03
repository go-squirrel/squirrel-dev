package infra

import "embed"

//go:embed scripts/*.sh
var scriptFS embed.FS
