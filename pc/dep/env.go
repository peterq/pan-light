package dep

import (
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type tEnv struct {
	Dev               bool   // 是否为开发环境
	Platform          string // 运行平台, darwin, windows, linux
	ApiBase           string // api调用前缀
	InternalServerUrl string
	VersionString     string
	DataPath          string
	Version           int
	ClientUA          string
}

func init() {
	rand.Seed(time.Now().UnixNano())
	Env.Platform = runtime.GOOS
	switch Env.Platform {
	case "linux", "darwin":
		p, _ := filepath.Abs(os.Getenv("HOME") + "/pt-program/pan-light")
		Env.DataPath = p
	case "windows":
		Env.DataPath = windowsHomeDir() + "/pt-program/pan-light"
	default:
		panic("unknown platform: " + Env.Platform)
	}
	Env.ClientUA += Env.Platform

	if exist, _ := pathExists(Env.DataPath); !exist {
		e := os.MkdirAll(Env.DataPath, os.ModePerm)
		if e != nil {
			Fatal(e.Error())
		}
	}
}

func DataPath(path string) string {
	return Env.DataPath + string(filepath.Separator) + strings.Join(strings.Split(path, "/"), string(filepath.Separator))
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func windowsHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
