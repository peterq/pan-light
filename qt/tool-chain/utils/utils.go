package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func ExistsFile(name string) bool {
	_, err := ioutil.ReadFile(name)
	return err == nil
}

func ExistsDir(name string) bool {
	_, err := ioutil.ReadDir(name)
	return err == nil
}

func MkdirAll(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		Log.WithError(err).Panicf("failed to create dir %v", dir)
	}
	return err
}

func RemoveAll(name string) error {
	err := os.RemoveAll(name)
	if err != nil {
		Log.WithError(err).Panicf("failed to remove %v", name)
	}
	return err
}

func Save(name, data string) error {
	err := ioutil.WriteFile(name, []byte(data), 0644)
	if err != nil {
		Log.WithError(err).Panicf("failed to save %v", name)
	} else {
		Log.Debugf("saved file len(%v) %v", len(data), name)
	}
	return err
}

func SaveExec(name, data string) error {
	err := ioutil.WriteFile(name, []byte(data), 0755)
	if err != nil {
		Log.WithError(err).Panicf("failed to save %v", name)
	} else {
		Log.Debugf("saved file len(%v) %v", len(data), name)
	}
	return err
}

func SaveBytes(name string, data []byte) error {
	err := ioutil.WriteFile(name, data, 0644)
	if err != nil {
		Log.WithError(err).Panicf("failed to save %v", name)
	}
	return err
}

//TODO: export error
func Load(name string) string {
	out, err := ioutil.ReadFile(name)
	if err != nil {
		Log.WithError(err).Errorf("failed to load %v", name)
		debug.PrintStack()
		os.Exit(0)
	}
	return string(out)
}

//TODO: export error
func LoadOptional(name string) string {
	out, err := ioutil.ReadFile(name)
	if err != nil {
		Log.WithError(err).Debugf("failed to load (optional) %v", name)
	}
	return string(out)
}

var (
	goQtPkgPath      string
	goQtPkgPathMutex = new(sync.Mutex)
)

func GoQtPkgPath(s ...string) (r string) {
	goQtPkgPathMutex.Lock()
	if len(goQtPkgPath) == 0 {
		goQtPkgPath = strings.TrimSpace(RunCmd(GoList("{{.Dir}}", packageName), "utils.GoQtPkgPath"))
		fmt.Println(goQtPkgPath)
		//os.Exit(0)
	}
	r = goQtPkgPath
	goQtPkgPathMutex.Unlock()
	return filepath.Join(r, filepath.Join(s...))
}

//TODO: export error
func RunCmd(cmd *exec.Cmd, name string) string {
	fields := logrus.Fields{"_func": "RunCmd", "name": name, "cmd": strings.Join(cmd.Args, " "), "env": strings.Join(cmd.Env, " "), "dir": cmd.Dir}
	Log.WithFields(fields).Debug("Execute")
	out, err := runCmdHelper(cmd)
	if err != nil {
		Log.WithError(err).WithFields(fields).Error("failed to run command")
		println(string(out))
		if ee, ok := err.(*exec.ExitError); ok {
			log.Println(string(ee.Stderr))
		}
		os.Exit(1)
	}
	return string(out)
}

//TODO: export error
func RunCmdOptional(cmd *exec.Cmd, name string) string {
	fields := logrus.Fields{"_func": "RunCmdOptional", "name": name, "cmd": strings.Join(cmd.Args, " "), "env": strings.Join(cmd.Env, " "), "dir": cmd.Dir}
	Log.WithFields(fields).Debug("Execute")
	out, err := runCmdHelper(cmd)
	if err != nil && !strings.Contains(string(out), "No template (-t) specified") {
		Log.WithError(err).WithFields(fields).Debug("failed to run command")
		if Log.Level == logrus.DebugLevel {
			println(string(out))
		}
	}
	return string(out)
}

func RunCmdOptionalError(cmd *exec.Cmd, name string) (string, error) {
	fields := logrus.Fields{"_func": "RunCmdOptionalError", "name": name, "cmd": strings.Join(cmd.Args, " "), "env": strings.Join(cmd.Env, " "), "dir": cmd.Dir}
	Log.WithFields(fields).Debug("Execute")
	out, err := runCmdHelper(cmd)
	if err != nil {
		Log.WithError(err).WithFields(fields).Debug("failed to run command")
		if Log.Level == logrus.DebugLevel {
			println(string(out))
		}
	}
	return string(out), err
}

func runCmdHelper(cmd *exec.Cmd) (out []byte, err error) {
	if _, ok := os.LookupEnv("WINEDEBUG"); ok {
		go func() { out, err = cmd.CombinedOutput() }()
		for range time.NewTicker(250 * time.Millisecond).C {
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				break
			}
		}
		return
	}
	return cmd.Output()
	//return cmd.CombinedOutput()
}
