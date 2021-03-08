package main

import (
	"fmt"
	"os"

	"github.com/mberwanger/demux-proxy/cmd"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	cmd.Execute(
		buildVersion(version, commit, date),
		os.Exit,
		os.Args[1:],
	)
}

func buildVersion(version, commit, date string) string {
	result := version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	return result
}
