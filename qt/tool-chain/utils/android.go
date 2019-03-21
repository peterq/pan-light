package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

//TODO: support for JDK 9
func JDK_DIR() string {
	if dir, ok := os.LookupEnv("JDK_DIR"); ok {
		return filepath.Clean(dir)
	}
	if dir, ok := os.LookupEnv("JAVA_HOME"); ok {
		return filepath.Clean(dir)
	}
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf(windowsSystemDrive()+"\\Program Files\\Java\\jdk%v", strings.Split(RunCmd(exec.Command("java", "-version"), "deploy.jdk"), "\"")[1])
	case "darwin":
		return fmt.Sprintf("/Library/Java/JavaVirtualMachines/jdk%v.jdk/Contents/Home", strings.Split(RunCmd(exec.Command("java", "-version"), "deploy.jdk"), "\"")[1])
	default:
		return filepath.Join(os.Getenv("HOME"), "jdk")
	}
}

func ANDROID_SDK_DIR() string {
	if dir, ok := os.LookupEnv("ANDROID_SDK_DIR"); ok {
		return filepath.Clean(dir)
	}
	if dir, ok := os.LookupEnv("ANDROID_SDK_ROOT"); ok {
		return filepath.Clean(dir)
	}
	switch runtime.GOOS {
	case "windows":
		return windowsSystemDrive() + "\\android-sdk-windows"
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "android-sdk-macosx")
	default:
		return filepath.Join(os.Getenv("HOME"), "android-sdk-linux")
	}
}

func ANDROID_NDK_DIR() string {
	if dir, ok := os.LookupEnv("ANDROID_NDK_DIR"); ok {
		return filepath.Clean(dir)
	}
	if dir, ok := os.LookupEnv("ANDROID_NDK_ROOT"); ok {
		return filepath.Clean(dir)
	}
	if runtime.GOOS == "windows" {
		return windowsSystemDrive() + "\\android-ndk-r18b"
	}
	return filepath.Join(os.Getenv("HOME"), "android-ndk-r18b")
}
