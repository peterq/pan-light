// +build !plugin

package main

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/gui"
)

//go:generate protoc --go_out=. storage/types.proto
//go:generate protoc --go_out=. downloader/internal/types.proto
func main() {
	defer func() {
		dep.DoClose()
	}()
	dep.DoInit()
	gui.StartGui()
}
