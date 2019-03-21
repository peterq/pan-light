package deploy

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func run(target, name, depPath, device string) {
	switch target {
	case "android", "android-emulator":
		if utils.ExistsFile(filepath.Join(depPath, "build-debug.apk")) {
			exec.Command(filepath.Join(utils.ANDROID_SDK_DIR(), "platform-tools", "adb"), "install", "-r", filepath.Join(depPath, "build-debug.apk")).Start()
		} else {
			exec.Command(filepath.Join(utils.ANDROID_SDK_DIR(), "platform-tools", "adb"), "install", "-r", filepath.Join(depPath, "build-release-signed.apk")).Start()
		}

		//TODO: parse manifest for ident and start app (+ logcat)

	case "ios-simulator":
		if device == "" {
			out, _ := exec.Command("xcrun", "instruments", "-s").Output()
			lines := strings.Split(string(out), "iPhone")
			device = strings.Split(strings.Split(string(out), "iPhone 8 ("+strings.Split(strings.Split(lines[len(lines)-1], "(")[1], ")")[0]+") [")[1], "]")[0]
		}
		go utils.RunCmdOptional(exec.Command("xcrun", "instruments", "-w", device), "start simulator")
		time.Sleep(1 * time.Second)
		utils.RunCmdOptional(exec.Command("xcrun", "simctl", "uninstall", "booted", filepath.Join(depPath, "main.app")), "uninstall old app")
		utils.RunCmdOptional(exec.Command("xcrun", "simctl", "install", "booted", filepath.Join(depPath, "main.app")), "install new app")
		utils.RunCmdOptional(exec.Command("xcrun", "simctl", "launch", "booted", strings.Replace(name, "_", "", -1)), "start app") //TODO: parse ident from plist

	case "darwin":
		exec.Command("open", filepath.Join(depPath, fmt.Sprintf("%v.app", name))).Start()

	case "linux":
		exec.Command(filepath.Join(depPath, name)).Start()

	case "windows":
		if runtime.GOOS == target {
			exec.Command(filepath.Join(depPath, name+".exe")).Start()
		} else {
			exec.Command("wine", filepath.Join(depPath, name+".exe")).Start()
		}

	case "sailfish-emulator":
		if utils.QT_SAILFISH() {
			return
		}
		utils.RunCmdOptional(exec.Command(filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "registervm", filepath.Join(utils.SAILFISH_DIR(), "emulator", "Sailfish OS Emulator", "Sailfish OS Emulator.vbox")), "register vm")
		utils.RunCmdOptional(exec.Command(filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "sharedfolder", "add", "Sailfish OS Emulator", "--name", "GOPATH", "--hostpath", utils.MustGoPath(), "--automount"), "mount GOPATH")

		if runtime.GOOS == "windows" {
			utils.RunCmdOptional(exec.Command(filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "startvm", "Sailfish OS Emulator"), "start emulator")
		} else {
			utils.RunCmdOptional(exec.Command("nohup", filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "startvm", "Sailfish OS Emulator"), "start emulator")
		}

		time.Sleep(10 * time.Second)

		err := sailfish_ssh("2223", "nemo", "sudo", "rpm", "-i", "--force", strings.Replace(strings.Replace(depPath, utils.MustGoPath(), "/media/sf_GOPATH/", -1)+"/*.rpm", "\\", "/", -1))
		if err != nil {
			utils.Log.WithError(err).Errorf("failed to install %v for %v", name, target)
		}

		err = sailfish_ssh("2223", "nemo", "nohup", "/usr/bin/harbour-"+name, ">", "/dev/null", "2>&1", "&")
		if err != nil {
			utils.Log.WithError(err).Errorf("failed to run %v for %v", name, target)
		}

	case "js", "wasm": //TODO: REVIEW and use emscripten wrapper instead
		if runtime.GOOS == "darwin" {
			exec.Command("/Applications/Firefox Nightly.app/Contents/MacOS/firefox", filepath.Join(depPath, "index.html")).Start()
		}
	}
}
