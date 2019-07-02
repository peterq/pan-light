// +build !plugin

package main

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/gui"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

//go:generate protoc --go_out=. storage/types.proto
//go:generate protoc --go_out=. downloader/internal/types.proto
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if os.Args[0] != startCmd {
		master()
	}
	log.Println("pan-light process")
	defer func() {
		dep.DoClose()
	}()
	dep.DoInit()
	gui.StartGui()
}

const startCmd = "pan_light_start"

func master() {

	if runtime.GOOS != "linux" {
		windosMaster()
		return
	}

	log.Println("master process")
START_PAN:
	c := exec.Command(os.Args[0], os.Args[1:]...)
	c.Args[0] = startCmd
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	err := c.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code := status.ExitStatus()
				if code == 2 {
					goto START_PAN
				}
			}
		}
	}
	log.Println(err)
	os.Exit(0)
}

func windosMaster() {
	log.Println("master process")
START_PAN:
	os.Remove(dep.DataPath("reboot"))
	c := exec.Command(os.Args[0], os.Args[1:]...)
	c.Args[0] = startCmd
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	err := c.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code := status.ExitStatus()
				if code == 3221225477 {
					if os.Getenv("pan_light_render_exception_fix") != "true" {
						os.Setenv("pan_light_render_exception_fix", "true")
						goto START_PAN
					}
				}
			}
		}
	}
	bin, _ := ioutil.ReadFile(dep.DataPath("reboot"))
	if string(bin) == "true" {
		goto START_PAN
	}
	log.Println(err)
	os.Exit(0)
}
