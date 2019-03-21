package setup

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func Update() {
	utils.Log.Info("running: 'qtsetup update'")

	utils.RunCmd(exec.Command("go", "clean", "-i", "github.com/peterq/pan-light/qt/cmd/..."), "run \"go clean cmd\"")
	utils.RunCmd(exec.Command("go", "clean", "-i", "github.com/peterq/pan-light/qt/tool-chain/..."), "run \"go clean tool-chain\"")

	fetch := exec.Command("git", "fetch", "-f", "--all")
	fetch.Dir = filepath.Join(utils.MustGoPath(), "src", "github.com", "therecipe", "qt")
	utils.RunCmd(fetch, "run \"git fetch\"")

	checkoutCmd := exec.Command("git", "checkout", "-f", "--", utils.GoQtPkgPath("cmd"))
	checkoutCmd.Dir = filepath.Join(utils.MustGoPath(), "src", "github.com", "therecipe", "qt")
	utils.RunCmd(checkoutCmd, "run \"git checkout cmd\"")

	checkoutInternal := exec.Command("git", "checkout", "-f", "--", utils.GoQtPkgPath("tool-chain"))
	checkoutInternal.Dir = filepath.Join(utils.MustGoPath(), "src", "github.com", "therecipe", "qt")
	utils.RunCmd(checkoutInternal, "run \"git checkout tool-chain\"")

	hash := "please install git"
	if _, err := exec.LookPath("git"); err == nil {
		cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
		cmd.Dir = utils.GoQtPkgPath()
		hash = strings.TrimSpace(utils.RunCmdOptional(cmd, "get git hash"))
	}

	utils.RunCmd(exec.Command("go", "install", "-v", fmt.Sprintf("-ldflags=\"-X=github.com/peterq/pan-light/qt/tool-chain/cmd.buildVersion=%v\"", hash), fmt.Sprintf("github.com/peterq/pan-light/qt/cmd/...")), "run \"go install\"")

	Prep()
}

func Upgrade() {
	utils.Log.Info("running: 'qtsetup upgrade'")

	utils.RunCmd(exec.Command("go", "clean", "-i", "github.com/peterq/pan-light/qt/..."), "run \"go clean\"")
	utils.RemoveAll(utils.GoQtPkgPath())

	utils.RunCmd(exec.Command("go", "get", "-v", "-d", "github.com/peterq/pan-light/qt/cmd/..."), "run \"go get\"")

	hash := "please install git"
	if _, err := exec.LookPath("git"); err == nil {
		cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
		cmd.Dir = utils.GoQtPkgPath()
		hash = strings.TrimSpace(utils.RunCmdOptional(cmd, "get git hash"))
	}

	utils.RunCmd(exec.Command("go", "install", "-v", fmt.Sprintf("-ldflags=\"-X=github.com/peterq/pan-light/qt/tool-chain/cmd.buildVersion=%v\"", hash), fmt.Sprintf("github.com/peterq/pan-light/qt/cmd/...")), "run \"go install\"")
}
