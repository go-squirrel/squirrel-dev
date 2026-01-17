package script

import "embed"

//go:embed scripts/*.sh
var ScriptFS embed.FS
