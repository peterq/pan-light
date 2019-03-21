package utils

import (
	"os"
	"path/filepath"
)

func QT_RPI() bool {
	return os.Getenv("QT_RPI") == "true"
}

func RPI_COMPILER() string {
	if name, ok := os.LookupEnv("RPI_COMPILER"); ok {
		return name
	}
	return "arm-rpi-4.9.3-linux-gnueabihf"
}

func RPI_TOOLS_DIR() string {
	if dir, ok := os.LookupEnv("RPI_TOOLS_DIR"); ok {
		return filepath.Clean(dir)
	}
	return filepath.Join(os.Getenv("HOME"), "raspi/tools")
}

func RPI1_SYSROOT_DIR() string {
	if dir, ok := os.LookupEnv("RPI1_SYSROOT_DIR"); ok {
		return filepath.Clean(dir)
	}
	return filepath.Join(os.Getenv("HOME"), "raspi/sysroot")
}

func RPI2_SYSROOT_DIR() string {
	if dir, ok := os.LookupEnv("RPI2_SYSROOT_DIR"); ok {
		return filepath.Clean(dir)
	}
	return filepath.Join(os.Getenv("HOME"), "raspi/sysroot")
}

func RPI3_SYSROOT_DIR() string {
	if dir, ok := os.LookupEnv("RPI3_SYSROOT_DIR"); ok {
		return filepath.Clean(dir)
	}
	return filepath.Join(os.Getenv("HOME"), "raspi/sysroot")
}
