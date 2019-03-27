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
	ApiHost           string
	InternalServerUrl string
	VersionString     string
	DataPath          string
	Version           int
	ClientUA          string
	ElectronSecretUA  string
	ListenPort        int
}

var Env = tEnv{
	Dev:               true,
	Platform:          "",
	ApiHost:           "http://127.0.0.1:9050",
	InternalServerUrl: "",
	VersionString:     "v1.0.0",
	DataPath:          "",
	Version:           20181113001,
	ClientUA:          "pan-light/v1.0.0;build 20181113001;",
	ElectronSecretUA:  "secret",
	ListenPort:        5678,
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
