package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const packageName = "github.com/peterq/pan-light/qt/bindings"

var mustGoPath string
var mustGoPathMutex = new(sync.Mutex)

//GOBIN returns the general GOBIN string
func GOBIN() string {
	if dir, ok := os.LookupEnv("GOBIN"); ok {
		return filepath.Clean(dir)
	}
	return filepath.Join(MustGoPath(), "bin")
}

// MustGoPath returns the GOPATH that holds this package
// it exits if any error occurres and also caches the result
func MustGoPath() string {
	mustGoPathMutex.Lock()
	if len(mustGoPath) == 0 {
		mustGoPath = strings.TrimSpace(RunCmd(GoList("{{.Root}}", "github.com/peterq/pan-light/qt"), "get list gopath"))
		if len(mustGoPath) == 0 {
			mustGoPath = GOPATH()
		}
	}
	mustGoPathMutex.Unlock()
	return mustGoPath
}

// GOPATH returns the general GOPATH string
func GOPATH() string {
	if dir, ok := os.LookupEnv("GOPATH"); ok {
		return dir
	}

	home := "HOME"
	if runtime.GOOS == "windows" {
		home = "USERPROFILE"
	}
	if dir, ok := os.LookupEnv(home); ok {
		return filepath.Join(dir, "go")
	}

	return ""
}
