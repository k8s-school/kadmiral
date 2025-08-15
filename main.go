package main

import (
	"embed"

	"github.com/example/kadmiral/cmd"
)

//go:embed resource/*
var ResourceFS embed.FS

func main() {
	cmd.Execute()
}
