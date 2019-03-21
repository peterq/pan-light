package setup

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
	"github.com/peterq/pan-light/qt/tool-chain/binding/templater"

	"github.com/peterq/pan-light/qt/tool-chain/cmd"
	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func Generate(target string, docker, vagrant bool) {
	utils.Log.Infof("running: 'qtsetup generate %v' [docker=%v] [vagrant=%v]", target, docker, vagrant)
	if docker {
		cmd.Docker([]string{"/home/user/work/bin/qtsetup", "-debug", "generate"}, "linux", "", true)
		return
	}

	parser.LoadModules()

	mode := "full"
	switch {
	case target == "js", target == "wasm":

	case target != runtime.GOOS:
		mode = "cgo"

	case utils.QT_STUB():
		mode = "stub"
	}

	if target == "windows" && runtime.GOOS == target {
		os.Setenv("QT_DEBUG_CONSOLE", "true")
	}

	for _, module := range parser.GetLibs() {
		if !parser.ShouldBuildForTarget(module, target) {
			utils.Log.Debugf("skipping generation of %v for %v", module, target)
			continue
		}

		var license string
		switch module {
		case "Charts", "DataVisualization":
			license = strings.Repeat(" ", 20-len(module)) + "[GPLv3]"
		}
		utils.Log.Infof("generating %v qt/%v %v", mode, strings.ToLower(module), license)

		if target == runtime.GOOS || utils.QT_FAT() || (mode == "full" && (target == "js" || target == "wasm")) { //TODO: REVIEW
			templater.GenModule(module, target, templater.NONE)
		} else {
			templater.CgoTemplate(module, "", target, templater.MINIMAL, "", "") //TODO: collect errors
		}

		if utils.QT_DYNAMIC_SETUP() && mode == "full" && (target != "js" && target != "wasm") {
			cc, _ := templater.ParseCgo(strings.ToLower(module), target)
			if cc != "" {
				cmd := exec.Command("go", "tool", "cgo", utils.GoQtPkgPath(strings.ToLower(module), strings.ToLower(module)+".go"))
				cmd.Dir = utils.GoQtPkgPath(strings.ToLower(module))
				utils.RunCmdOptional(cmd, fmt.Sprintf("failed to run cgo for %v (%v) on %v", target, strings.ToLower(module), runtime.GOOS))
			}
		}
	}
}
