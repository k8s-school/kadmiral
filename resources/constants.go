package resources

import (
	"embed"
	_ "embed"
)

//go:embed *
var Fs embed.FS
