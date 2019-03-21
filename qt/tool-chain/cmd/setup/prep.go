package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func Prep() {
	utils.Log.Info("running: 'qtsetup prep'")

	errString := "failed to create %v symlink in your PATH (%v); please use %v instead"
	sucString := "successfully created %v symlink in your PATH (%v)"

	for _, app := range []string{"qtrcc", "qtmoc", "qtminimal", "qtdeploy", "go"} {
		if app == "go" && !(utils.QT_MSYS2() && !utils.QT_DOCKER()) {
			continue
		}

		if runtime.GOOS == "windows" {
			sPath := filepath.Join(utils.GOBIN(), fmt.Sprintf("%v.exe", app))
			dPath := filepath.Join(runtime.GOROOT(), "bin", fmt.Sprintf("%v.exe", app))
			if utils.QT_MSYS2() && !utils.QT_DOCKER() {
				if app == "go" {
					sPath = dPath
				}
				dPath = filepath.Join(utils.QT_MSYS2_DIR(), "..", "usr", "bin", fmt.Sprintf("%v.exe", app))
			}
			if sPath == dPath {
				continue
			}
			utils.RemoveAll(dPath)
			if err := os.Link(sPath, dPath); err == nil {
				utils.Log.Infof(sucString, app, dPath)
			} else {
				utils.Log.Warnf(errString, app, dPath, sPath)
			}
			continue
		} else {
			var suc bool
			sPath := filepath.Join(utils.GOBIN(), app)
			var dPath string
			for _, pdPath := range filepath.SplitList("/usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin:" + filepath.Join(filepath.Join(runtime.GOROOT(), "bin"))) {
				dPath = filepath.Join(pdPath, app)
				if sPath == dPath {
					continue
				}
				utils.RemoveAll(dPath)
				if err := os.Symlink(sPath, dPath); err == nil {
					suc = true
					break
				}
			}
			if suc {
				utils.Log.Infof(sucString, app, dPath)
			} else {
				utils.Log.Warnf(errString, app, dPath, sPath)
			}
		}
	}
}
