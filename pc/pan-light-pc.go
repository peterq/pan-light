// +build !plugin

package main

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/gui"
)

func main() {
	defer func() {
		dep.DoClose()
	}()
	dep.DoInit()
	gui.StartGui()
}
