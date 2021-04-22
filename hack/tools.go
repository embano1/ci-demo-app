// +build tools

package tools

// This package imports things required by this repository, to force `go mod` to see them as dependencies
import (
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
)
