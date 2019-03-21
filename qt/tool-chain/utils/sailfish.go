package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func VIRTUALBOX_DIR() string {
	if dir, ok := os.LookupEnv("VIRTUALBOX_DIR"); ok {
		return filepath.Clean(dir)
	}
	if runtime.GOOS == "windows" {
		return windowsSystemDrive() + "\\Program Files\\Oracle\\VirtualBox"
	}
	path, err := exec.LookPath("vboxmanage")
	if err != nil {
		Log.WithError(err).Error("failed to find vboxmanage in your PATH")
	}
	path = filepath.Dir(path)
	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			Log.WithError(err).WithField("path", path).Fatal("can't resolve absolute path")
		}
	}
	return path
}

func SAILFISH_DIR() string {
	if dir, ok := os.LookupEnv("SAILFISH_DIR"); ok {
		return filepath.Clean(dir)
	}
	if runtime.GOOS == "windows" {
		return windowsSystemDrive() + "\\SailfishOS"
	}
	return filepath.Join(os.Getenv("HOME"), "SailfishOS")
}

func QT_SAILFISH() bool {
	return os.Getenv("QT_SAILFISH") == "true"
}

func QT_SAILFISH_VERSION() string {
	if ver, ok := os.LookupEnv("QT_SAILFISH_VERSION"); ok {
		return ver
	}
	return "2.2.1.18"
}
