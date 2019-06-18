package dep

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	Env.Platform = runtime.GOOS
	switch Env.Platform {
	case "linux", "darwin":
		p, _ := filepath.Abs(os.Getenv("HOME") + "/pt-program/pan-light")
		Env.DataPath = p
	default:

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
