package resources

import (
	_ "embed"
)

//go:embed init.sh
var InitScript string

//go:embed install-cilium.sh
var InstallCiliumScript string

//go:embed reset.sh
var ResetScript string
