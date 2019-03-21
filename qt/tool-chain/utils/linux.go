package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func QT_PKG_CONFIG() bool {
	return os.Getenv("QT_PKG_CONFIG") == "true"
}

func QT_DOC_DIR() string {
	if dir, ok := os.LookupEnv("QT_DOC_DIR"); ok {
		return filepath.Clean(dir)
	}
	switch QT_DISTRO() {
	case "arch":
		return "/usr/share/doc/qt"
	case "fedora":
		return "/usr/share/doc/qt5"
	case "suse":
		return "/usr/share/doc/packages/qt5"
	case "ubuntu":
		return "/usr/share/qt5/doc"
	case "gentoo":
		return "/usr/share/doc/qt-" + QT_VERSION()
	default:
		Log.Error("failed to detect the Linux distro")
		return ""
	}
}

func QT_MISC_DIR() string {
	if dir, ok := os.LookupEnv("QT_MISC_DIR"); ok {
		return filepath.Clean(dir)
	}
	if QT_DISTRO() == "arch" {
		return filepath.Join(strings.TrimSpace(RunCmd(exec.Command("pkg-config", "--variable=libdir", "Qt5Core"), "cgo.LinuxPkgConfig_libDir")), "qt")
	}
	//fedora, suse, ubuntu, gentoo
	return strings.TrimSuffix(strings.TrimSpace(RunCmd(exec.Command("pkg-config", "--variable=host_bins", "Qt5Core"), "cgo.LinuxPkgConfig_hostBins")), "/bin")
}

func QT_DISTRO() string {
	if distro, ok := os.LookupEnv("QT_DISTRO"); ok {
		return distro
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		return "arch"
	}
	if _, err := exec.LookPath("yum"); err == nil {
		return "fedora"
	}
	if _, err := exec.LookPath("zypper"); err == nil {
		return "suse"
	}
	if _, err := exec.LookPath("apt-get"); err == nil {
		return "ubuntu"
	}
	if _, err := exec.LookPath("emerge"); err == nil {
		return "gentoo"
	}
	Log.Error("failed to detect the Linux distro")
	return ""
}

func QT_MXE_ARCH() string {
	if arch := os.Getenv("QT_MXE_ARCH"); arch == "amd64" {
		return arch
	}
	return "386"
}

func QT_MXE_STATIC() bool {
	return os.Getenv("QT_MXE_STATIC") == "true"
}

func QT_MXE_TRIPLET() string {
	prefix := "i686"
	if QT_MXE_ARCH() == "amd64" {
		prefix = "x86_64"
	}
	suffix := "shared"
	if QT_MXE_STATIC() {
		suffix = "static"
	}
	return fmt.Sprintf("%v-w64-mingw32.%v", prefix, suffix)
}

func QT_MXE_DIR() string {
	if dir, ok := os.LookupEnv("QT_MXE_DIR"); ok {
		return filepath.Clean(dir)
	}
	return filepath.Join("/usr", "lib", "mxe")
}

func QT_MXE_BIN(tool string) string {
	return filepath.Join(QT_MXE_DIR(), "usr", "bin", fmt.Sprintf("%v-%v", QT_MXE_TRIPLET(), tool))
}

func QT_MXE() bool {
	return os.Getenv("QT_MXE") == "true"
}

func QT_UBPORTS() bool {
	return os.Getenv("QT_UBPORTS") == "true"
}

func QT_UBPORTS_ARCH() string {
	if arch := os.Getenv("QT_UBPORTS_ARCH"); arch == "amd64" {
		return arch
	}
	return "arm"
}

func QT_UBPORTS_VERSION() string {
	if rel := os.Getenv("QT_UBPORTS_VERSION"); rel == "xenial" {
		return rel
	}
	return "vivid"
}
