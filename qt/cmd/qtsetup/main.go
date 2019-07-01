package main

import (
	"flag"
	"fmt"
	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
	"os"
	"runtime"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/cmd"
	"github.com/peterq/pan-light/qt/tool-chain/cmd/setup"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func main() {
	flag.Usage = func() {
		println("Usage: qtsetup [-debug] [mode] [target]\n")

		println("Flags:\n")
		flag.PrintDefaults()
		println()

		println("Modes:\n")
		for _, m := range []struct{ name, desc string }{
			{"prep", "symlink tooling into the PATH"},
			{"check", "perform some basic env checks"},
			{"generate", "generate code for all packages"},
			{"install", "go install all packages"},
			{"test", "build and test some examples"},
			{"full", "run all of the above"},
			{"help", "print help"},
			{"update", "update 'cmd' and 'tool-chain/cmd'"},
			{"upgrade", "update everything"},
		} {
			fmt.Printf("  %v%v%v\n", m.name, strings.Repeat(" ", 12-len(m.name)), m.desc)
		}
		println()

		println("Targets:\n")
		//TODO:
		println()

		os.Exit(0)
	}

	var docker bool
	flag.BoolVar(&docker, "docker", false, "run command inside docker container")

	var vagrant bool
	flag.BoolVar(&vagrant, "vagrant", false, "run command inside vagrant vm")

	var dynamic bool
	if runtime.GOOS != "windows" {
		flag.BoolVar(&dynamic, "dynamic", false, "create and use semi-dynamic libraries during the generation and installation process (experimental; no real replacement for dynamic linking)")
	}

	if cmd.ParseFlags() {
		flag.Usage()
	}

	mode := "full"
	target := runtime.GOOS

	switch flag.NArg() {
	case 0:
	case 1:
		mode = flag.Arg(0)
	case 2:
		mode = flag.Arg(0)
		target = flag.Arg(1)
	default:
		flag.Usage()
	}

	var vagrant_system string
	if target_splitted := strings.Split(target, "/"); vagrant && len(target_splitted) == 2 {
		vagrant_system = target_splitted[0]
		target = target_splitted[1]
	}

	if target == "desktop" {
		target = runtime.GOOS
	}
	utils.CheckBuildTarget(target)
	cmd.InitEnv(target)

	if dynamic && (target == runtime.GOOS || target == "js" || target == "wasm") {
		os.Setenv("QT_DYNAMIC_SETUP", "true")
	}

	if target == "js" || target == "wasm" { //TODO: remove for module support + resolve dependencies
		os.Setenv("GOCACHE", "off")
	}

	parser.State.Target = target

	switch mode {
	case "prep":
		setup.Prep()
	case "check":
		setup.Check(target, docker, vagrant)
	case "generate":
		setup.Generate(target, docker, vagrant)
	case "install":
		setup.Install(target, docker, vagrant)
	case "test":
		setup.Test(target, docker, vagrant, vagrant_system)
	case "full":
		setup.Prep()
		setup.Check(target, docker, vagrant)
		setup.Generate(target, docker, vagrant)
		setup.Install(target, docker, vagrant)
		setup.Test(target, docker, vagrant, vagrant_system)
	case "update":
		setup.Update()
	case "upgrade":
		setup.Upgrade()
	default:
		flag.Usage()
	}
}
