package deploy

import (
	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/peterq/pan-light/qt/tool-chain/cmd"
	"github.com/peterq/pan-light/qt/tool-chain/cmd/minimal"
	"github.com/peterq/pan-light/qt/tool-chain/cmd/moc"
	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func Deploy(mode, target, path string, docker bool, ldFlags, tags string, fast bool, device string, vagrant bool, vagrantsystem string, comply bool) {
	parser.State.Target = target
	utils.Log.WithField("mode", mode).WithField("target", target).WithField("path", path).WithField("docker", docker).WithField("ldFlags", ldFlags).WithField("fast", fast).WithField("comply", comply).Debug("running Deploy")
	name := filepath.Base(path)
	switch name {
	case "lib", "plugins", "qml",
		"audio", "bearer", "iconengines", "imageformats", "mediaservice",
		"platforminputcontexts", "platforms", "playlistformats", "qmltooling",
		"qt", "Qt", "QT", "styles", "translations":
		name += "_project"
	}
	depPath := filepath.Join(path, "deploy", target)

	switch mode {
	case "build", "test":

		if docker || vagrant {
			args := []string{"qtdeploy", "-debug"}
			if fast {
				args = append(args, "-fast")
			}
			if comply {
				args = append(args, "-comply")
			}
			if vagrantsystem == "docker" {
				args = append(args, "-docker")
			}
			args = append(args, []string{"-ldflags=" + ldFlags, "-tags=" + tags, "build"}...)

			if docker {
				cmd.Docker(args, target, path, false)
			} else {
				cmd.Vagrant(args, target, path, false, vagrantsystem)
			}
			break
		}

		if !fast {
			err := os.RemoveAll(depPath)
			if err != nil {
				utils.Log.WithError(err).Panic("failed to remove deploy folder")
			}

			if utils.UseGOMOD(path) {
				if !utils.ExistsDir(filepath.Join(path, "vendor")) {
					cmd := exec.Command("go", "mod", "vendor")
					cmd.Dir = path
					utils.RunCmd(cmd, "go mod vendor")
				}
			}
		}

		if utils.ExistsDir(depPath + "_obj") {
			utils.RemoveAll(depPath + "_obj")
		}

		//rcc.Rcc(path, target, tags, os.Getenv("QTRCC_OUTPUT_DIR"))
		if false && !fast {
			moc.Moc(path, target, tags, false, false)
		}

		if false && ((!fast || utils.QT_STUB()) || ((target == "js" || target == "wasm") && (utils.QT_DOCKER() || utils.QT_VAGRANT()))) && !utils.QT_FAT() {
			minimal.Minimal(path, target, tags)
		}

		build(mode, target, path, ldFlags, tags, name, depPath, fast, comply)

		if !(fast || (utils.QT_DEBUG_QML() && target == runtime.GOOS)) || (target == "js" || target == "wasm") || true {
			bundle(mode, target, path, name, depPath, tags, fast)
		} else if fast {
			switch target {
			case "darwin":
				if fn := filepath.Join(depPath, name+".app", "Contents", "Info.plist"); !utils.ExistsFile(fn) {
					utils.Save(fn, darwin_plist(name))
				}
			}
		}
	}

	if (mode == "run" || mode == "test") && !(fast && (target == "js" || target == "wasm")) {
		run(target, name, depPath, device)
	}
}
