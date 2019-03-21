package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var sdkMutex = new(sync.Mutex)

func XCODE_DIR() string {
	if dir, ok := os.LookupEnv("XCODE_DIR"); ok {
		return filepath.Clean(dir)
	}
	return filepath.Join("/Applications/Xcode.app")
}

var _MACOS_SDK_DIR string

func MACOS_SDK_DIR() string {
	sdkMutex.Lock()
	defer sdkMutex.Unlock()
	if _MACOS_SDK_DIR != "" {
		return _MACOS_SDK_DIR
	}
	if runtime.GOOS == "darwin" {
		basePath := filepath.Join(XCODE_DIR(), "Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs")
		for maj := 10; maj < 50; maj++ {
			for min := 0; min < 25; min++ {
				i := fmt.Sprintf("%v.%v", maj, min)
				if ExistsDir(filepath.Join(basePath, fmt.Sprintf("MacOSX%v.sdk", i))) {
					_MACOS_SDK_DIR = fmt.Sprintf("MacOSX%v.sdk", i)
					return _MACOS_SDK_DIR
				}
			}
		}
		if ExistsDir(filepath.Join(basePath, "MacOSX.sdk")) {
			_MACOS_SDK_DIR = "MacOSX.sdk"
			return _MACOS_SDK_DIR
		}
		Log.Errorf("failed to find MacOSX sdk in %v", basePath)
	}
	return ""
}

var _IPHONEOS_SDK_DIR string

func IPHONEOS_SDK_DIR() string {
	sdkMutex.Lock()
	defer sdkMutex.Unlock()
	if _IPHONEOS_SDK_DIR != "" {
		return _IPHONEOS_SDK_DIR
	}
	if runtime.GOOS == "darwin" {
		basePath := filepath.Join(XCODE_DIR(), "Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs")
		for maj := 10; maj < 50; maj++ {
			for min := 0; min < 25; min++ {
				i := fmt.Sprintf("%v.%v", maj, min)
				if ExistsDir(filepath.Join(basePath, fmt.Sprintf("iPhoneOS%v.sdk", i))) {
					_IPHONEOS_SDK_DIR = fmt.Sprintf("iPhoneOS%v.sdk", i)
					return _IPHONEOS_SDK_DIR
				}
			}
		}
		if ExistsDir(filepath.Join(basePath, "iPhoneOS.sdk")) {
			_IPHONEOS_SDK_DIR = "iPhoneOS.sdk"
			return _IPHONEOS_SDK_DIR
		}
		Log.Errorf("failed to find iPhoneOS sdk in %v", basePath)
	}
	return ""
}

var _IPHONESIMULATOR_SDK_DIR string

func IPHONESIMULATOR_SDK_DIR() string {
	sdkMutex.Lock()
	defer sdkMutex.Unlock()
	if _IPHONESIMULATOR_SDK_DIR != "" {
		return _IPHONESIMULATOR_SDK_DIR
	}
	if runtime.GOOS == "darwin" {
		basePath := filepath.Join(XCODE_DIR(), "Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs")
		for maj := 10; maj < 50; maj++ {
			for min := 0; min < 25; min++ {
				i := fmt.Sprintf("%v.%v", maj, min)
				if ExistsDir(filepath.Join(basePath, fmt.Sprintf("iPhoneSimulator%v.sdk", i))) {
					_IPHONESIMULATOR_SDK_DIR = fmt.Sprintf("iPhoneSimulator%v.sdk", i)
					return _IPHONESIMULATOR_SDK_DIR
				}
			}
		}
		if ExistsDir(filepath.Join(basePath, "iPhoneSimulator.sdk")) {
			_IPHONESIMULATOR_SDK_DIR = "iPhoneSimulator.sdk"
			return _IPHONESIMULATOR_SDK_DIR
		}
		Log.Errorf("failed to find iPhoneSimulator sdk in %v", basePath)
	}
	return ""
}

func QT_HOMEBREW() bool {
	return os.Getenv("QT_HOMEBREW") == "true" || isHomeBrewQtDir()
}

func QT_MACPORTS() bool {
	return os.Getenv("QT_MACPORTS") == "true"
}

func QT_NIX() bool {
	_, ok := os.LookupEnv("NIX_STORE")
	return ok
}

func isHomeBrewQtDir() bool {
	return ExistsFile(filepath.Join(QT_DIR(), "INSTALL_RECEIPT.json"))
}

func QT_DARWIN_DIR() string {
	path := qT_DARWIN_DIR()
	if ExistsDir(path) {
		return path
	}
	return strings.Replace(path, QT_VERSION_MAJOR(), QT_VERSION(), -1)
}

var qt_darwin_dir_nix string

func qT_DARWIN_DIR() string {
	if QT_HOMEBREW() {
		if isHomeBrewQtDir() {
			return QT_DIR()
		}
		return "/usr/local/opt/qt5"
	}
	if QT_MACPORTS() {
		return "/opt/local/libexec/qt5"
	}
	if QT_NIX() {
		if len(qt_darwin_dir_nix) == 0 {
			qt_darwin_dir_nix = strings.TrimSpace(RunCmd(exec.Command(ToolPath("qmake", "darwin"), "-query", "QT_INSTALL_PREFIX"), "nix qt dir"))
		}
		return qt_darwin_dir_nix
	}
	return filepath.Join(QT_DIR(), fmt.Sprintf("%v/clang_64", QT_VERSION_MAJOR()))
}
