package deploy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/peterq/pan-light/qt/tool-chain/cmd"
	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func build(mode, target, path, ldFlagsCustom, tagsCustom, name, depPath string, fast, comply bool) {
	env, tags, ldFlags, out := cmd.BuildEnv(target, name, depPath)
	if ((!fast || utils.QT_STUB()) && !utils.QT_FAT()) || target == "js" || target == "wasm" {
		tags = append(tags, "minimal")
	}
	if tagsCustom != "" {
		tags = append(tags, strings.Split(tagsCustom, " ")...)
	}
	if utils.QT_DEBUG_QML() && target == runtime.GOOS {
		out = filepath.Join(depPath, name)
	}

	var ending string
	switch target {
	case "android", "android-emulator", "ios", "ios-simulator":
		utils.Save(filepath.Join(path, "cgo_main_wrapper.go"), "package main\nimport (\n\"C\"\n\"os\"\n\"unsafe\"\n)\n//export go_main_wrapper\nfunc go_main_wrapper(argc C.int, argv unsafe.Pointer) {\nos.Args=make([]string,int(argc))\nfor i,b := range (*[1<<3]*C.char)(argv)[:int(argc):int(argc)] {\nos.Args[i] = C.GoString(b)\n}\nmain()\n}")
	case "windows":
		ending = ".exe"
	case "sailfish", "sailfish-emulator":
		if !utils.QT_SAILFISH() {
			build_sailfish(target, path, ldFlagsCustom, name)
			return
		}
	case "js":
		build_js(target, path, env, tags, out)
		return
	case "wasm":
		ending = ".wasm"
	case "linux":
		if fast || utils.QT_PKG_CONFIG() {
			delete(env, "CGO_LDFLAGS")
		}
	}

	var pattern string
	if strings.Contains(runtime.Version(), "1.1") || strings.Contains(runtime.Version(), "devel") {
		pattern = "all="
	}

	if utils.Log.Level == logrus.DebugLevel && target != "wasm" {
		ldFlags = append(ldFlags, "-extldflags=-v")
	}

	cmd := exec.Command("go", "build", "-p", strconv.Itoa(runtime.GOMAXPROCS(0)), "-v")
	if len(ldFlags) > 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("-ldflags=%v%v", pattern, escapeFlags(ldFlags, ldFlagsCustom)))
	}
	cmd.Args = append(cmd.Args, "-o", out+ending)

	cmd.Dir = path

	if fast && !utils.QT_STUB() {
		cmd.Args = append(cmd.Args, "-i")
	}

	if comply {
		utils.MkdirAll(depPath + "_obj")
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOTMPDIR=%v", depPath+"_obj"))
		cmd.Args = append(cmd.Args, "-a", "-x", "-work")
	} else if utils.Log.Level == logrus.DebugLevel {
		cmd.Args = append(cmd.Args, "-x")
	}

	cmd.Args = append(cmd.Args, fmt.Sprintf("-tags=\"%v\"", strings.Join(tags, "\" \"")))

	if target != runtime.GOOS {
		//cmd.Args = append(cmd.Args, []string{"-pkgdir", filepath.Join(utils.MustGoPath(), "pkg", fmt.Sprintf("%v_%v_%v", strings.Replace(target, "-", "_", -1), env["GOOS"], env["GOARCH"]))}...)
	}

	switch target {
	case "android", "android-emulator":
		cmd.Args = append(cmd.Args, "-buildmode", "c-shared")
	case "ios", "ios-simulator":
		cmd.Args = append(cmd.Args, "-buildmode", "c-archive")
	}

	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", key, value))
	}
	log.Println(cmd.Args)

	utils.RunCmd(cmd, fmt.Sprintf("build for %v on %v", target, runtime.GOOS))

	if target == "darwin" && !fast {
		strip := exec.Command("strip", "-x", out) //TODO: -u -r
		strip.Dir = path
		utils.RunCmd(strip, fmt.Sprintf("strip binary for %v on %v", target, runtime.GOOS))
	}

	utils.RemoveAll(filepath.Join(path, "cgo_main_wrapper.go"))

	if comply {
		dirs, err := ioutil.ReadDir(depPath + "_obj")
		if err != nil {
			utils.Log.WithError(err).Error("failed to read object dir")
		}

		var randname string
		for _, dir := range dirs {
			if strings.HasPrefix(dir.Name(), "go-build") {
				randname = dir.Name()
				os.Rename(filepath.Join(depPath+"_obj", dir.Name()), depPath+"_objreal")
				utils.RemoveAll(depPath + "_obj")
				os.Rename(depPath+"_objreal", depPath+"_obj")
				break
			}
		}

		walkFn := func(fpath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			switch info.Name() {
			case "_pkg_.a":
			case "importcfg.link":
				pre := utils.Load(fpath)
				pre = strings.Replace(pre, filepath.Join(depPath+"_obj", randname), ".", -1)
				utils.Save(fpath, pre)
			default:
				if !info.IsDir() || info.Name() == "exe" {
					utils.RemoveAll(fpath)
				}
			}
			return nil
		}
		filepath.Walk(depPath+"_obj", walkFn)

		utils.SaveExec(filepath.Join(depPath+"_obj", "relink.sh"), relink(env, target))
	}
}

func build_sailfish(target, path, ldFlagsCustom, name string) {
	//TODO: ldFlagsCustom, tags

	if !strings.Contains(path, utils.MustGoPath()) {
		utils.Log.Panicln("Project needs to be inside GOPATH; have:", path, "want:", utils.MustGoPath())
	}

	utils.RunCmdOptional(exec.Command(filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "registervm", filepath.Join(utils.SAILFISH_DIR(), "mersdk", "Sailfish OS Build Engine", "Sailfish OS Build Engine.vbox")), fmt.Sprintf("register mersdk for %v on %v", target, runtime.GOOS))
	utils.RunCmdOptional(exec.Command(filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "sharedfolder", "add", "Sailfish OS Build Engine", "--name", "GOROOT", "--hostpath", runtime.GOROOT(), "--automount"), fmt.Sprintf("share GOROOT dir for %v on %v", target, runtime.GOOS))
	utils.RunCmdOptional(exec.Command(filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "sharedfolder", "add", "Sailfish OS Build Engine", "--name", "GOPATH", "--hostpath", utils.MustGoPath(), "--automount"), fmt.Sprintf("share GOPATH dir for %v on %v", target, runtime.GOOS))

	if runtime.GOOS == "windows" {
		utils.RunCmdOptional(exec.Command(filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "startvm", "--type", "headless", "Sailfish OS Build Engine"), fmt.Sprintf("start vbox mersdk for %v on %v", target, runtime.GOOS))
	} else {
		utils.RunCmdOptional(exec.Command("nohup", filepath.Join(utils.VIRTUALBOX_DIR(), "vboxmanage"), "startvm", "--type", "headless", "Sailfish OS Build Engine"), fmt.Sprintf("start vbox mersdk for %v on %v", target, runtime.GOOS))
	}

	time.Sleep(10 * time.Second)

	for _, l := range []string{"libmpc.so.3", "libmpfr.so.4", "libgmp.so.10", "libpthread_nonshared.a", "libc_nonshared.a"} {
		sailfish_ssh("2222", "root", "ln", "-s", fmt.Sprintf("/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/usr/lib/%v", l), fmt.Sprintf("/usr/lib/%v", l))
	}

	arch, gcc := "i486", "gnu"
	if target == "sailfish" {
		arch, gcc = "armv7hl", "gnueabi"
	}

	sailfish_ssh("2222", "root", "ln", "-s", fmt.Sprintf("/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/opt/cross/bin/%v-meego-linux-%v-as", arch, gcc), fmt.Sprintf("/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/opt/cross/libexec/gcc/%v-meego-linux-%v/4.8.3/as", arch, gcc))
	sailfish_ssh("2222", "root", "ln", "-s", fmt.Sprintf("/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/opt/cross/bin/%v-meego-linux-%v-ld", arch, gcc), fmt.Sprintf("/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/opt/cross/libexec/gcc/%v-meego-linux-%v/4.8.3/ld", arch, gcc))

	var pattern string
	if strings.Contains(runtime.Version(), "1.1") || strings.Contains(runtime.Version(), "devel") {
		pattern = "all="
	}

	//TODO:
	var err error
	if target == "sailfish-emulator" {
		err = sailfish_ssh("2222", "root", "cd", strings.Replace(strings.Replace(path, utils.MustGoPath(), "/media/sf_GOPATH/", -1), "\\", "/", -1), "&&", "GOROOT=/media/sf_GOROOT", "GOPATH=/media/sf_GOPATH", "PATH=$PATH:$GOROOT/bin/linux_386", "GOOS=linux", "GOARCH=386", "CGO_ENABLED=1", "CGO_CFLAGS_ALLOW=.*", "CGO_CXXFLAGS_ALLOW=.*", "CGO_LDFLAGS_ALLOW=.*", "CC=/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/opt/cross/bin/i486-meego-linux-gnu-gcc", "CXX=/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/opt/cross/bin/i486-meego-linux-gnu-g++", "CPATH=/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/include", "LIBRARY_PATH=/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/lib:/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/lib:/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/lib/pulseaudio", "CGO_LDFLAGS=--sysroot=/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/", "go", "build", fmt.Sprintf("-ldflags=%v\"-s -w\"", pattern), "-tags=\"minimal sailfish_emulator\"", "-o", "deploy/"+target+"/harbour-"+name)
	} else {
		err = sailfish_ssh("2222", "root", "cd", strings.Replace(strings.Replace(path, utils.MustGoPath(), "/media/sf_GOPATH/", -1), "\\", "/", -1), "&&", "GOROOT=/media/sf_GOROOT", "GOPATH=/media/sf_GOPATH", "PATH=$PATH:$GOROOT/bin/linux_386", "GOOS=linux", "GOARCH=arm", "GOARM=7", "CGO_ENABLED=1", "CGO_CFLAGS_ALLOW=.*", "CGO_CXXFLAGS_ALLOW=.*", "CGO_LDFLAGS_ALLOW=.*", "CC=/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/opt/cross/bin/armv7hl-meego-linux-gnueabi-gcc", "CXX=/srv/mer/toolings/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"/opt/cross/bin/armv7hl-meego-linux-gnueabi-g++", "CPATH=/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-armv7hl/usr/include", "LIBRARY_PATH=/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-armv7hl/usr/lib:/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-armv7hl/lib:/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-armv7hl/usr/lib/pulseaudio", "CGO_LDFLAGS=--sysroot=/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-armv7hl/", "go", "build", "-x", "-v", fmt.Sprintf("-ldflags=%v\"-s -w\"", pattern), "-tags=\"minimal sailfish\"", "-o", "deploy/"+target+"/harbour-"+name)
	}
	if err != nil {
		println(err.Error())
		utils.Log.Panicf("failed to build for %v on %v", target, runtime.GOOS)
	}
}

func build_js(target string, path string, env map[string]string, tags []string, out string) {
	cmd := exec.Command(filepath.Join(utils.GOBIN(), "gopherjs"), "build", ".", "-v", "-m", "-o", filepath.Join(filepath.Dir(out), "go.js"))
	cmd.Dir = path

	//TODO (bug in gopherjs?): cmd.Args = append(cmd.Args, fmt.Sprintf("--tags=\"%v\"", strings.Join(tags[1:], " ")))
	cmd.Args = append(cmd.Args, fmt.Sprintf("--tags=%v", tags[1]))

	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", key, value))
	}

	utils.RunCmd(cmd, fmt.Sprintf("build for %v on %v", target, runtime.GOOS))
}
