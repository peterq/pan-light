package templater

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func cleanLibs(module string, mode int) []string {

	var out []string

	switch {
	case mode == RCC:
		out = []string{"Core"}
	case mode == MOC, module == "build_static":
		out = parser.LibDeps[module]
	case mode == MINIMAL, mode == NONE:
		out = append([]string{module}, parser.LibDeps[module]...)
	}

	for i, v := range out {
		if v == "Speech" {
			out[i] = "TextToSpeech"
		}
	}
	return out
}

//needed for static linking
func GetiOSClang(buildTarget, _, path string) []string {
	var tmp = CgoTemplate("build_static", path, buildTarget, NONE, "main", "")

	tmp = strings.Split(tmp, "/*")[1]
	tmp = strings.Split(tmp, "*/")[0]

	tmp = strings.Replace(tmp, "#cgo CFLAGS: ", "", -1)
	tmp = strings.Replace(tmp, "#cgo CXXFLAGS: ", "", -1)
	tmp = strings.Replace(tmp, "#cgo LDFLAGS: ", "", -1)
	tmp = strings.Replace(tmp, "\n", " ", -1)

	if buildTarget == "ios" {
		tmp = strings.Replace(tmp, "_iphonesimulator", "", -1)
		tmp = strings.Replace(tmp, "x86_64", "arm64", -1)
		tmp = strings.Replace(tmp, "iPhoneSimulator", "iPhoneOS", -1)
		tmp = strings.Replace(tmp, "ios-simulator", "iphoneos", -1)
	}

	return strings.Split(tmp, " ")
}

func cgoSailfish(module, mocPath string, mode int, pkg string, libs []string) {
	var bb = new(bytes.Buffer)
	defer bb.Reset()

	if mode != MOC {
		libs = cleanLibs(module, mode)
	}

	fmt.Fprintf(bb, "// +build ${BUILDTARGET}%v\n\n", func() string {
		if mode == MINIMAL {
			return ",minimal"
		}
		if mode == NONE {
			return ",!minimal"
		}
		return ""
	}())

	fmt.Fprintf(bb, "package %v\n\n", func() string {
		if mode == MOC {
			return pkg
		}
		return strings.ToLower(module)
	}())
	fmt.Fprint(bb, "/*\n")

	fmt.Fprint(bb, "#cgo CFLAGS: -pipe -O2 -g -pipe -Wall -Wp,-D_FORTIFY_SOURCE=2 -fexceptions -fstack-protector --param=ssp-buffer-size=4 -Wformat -Wformat-security -m32 -msse -msse2 -march=i686 -mfpmath=sse -mtune=generic -fno-omit-frame-pointer -fasynchronous-unwind-tables -fPIC -fvisibility=hidden -fvisibility-inlines-hidden -Wall -W -D_REENTRANT -fPIC\n")
	fmt.Fprint(bb, "#cgo CXXFLAGS: -pipe -O2 -g -pipe -Wall -Wp,-D_FORTIFY_SOURCE=2 -fexceptions -fstack-protector --param=ssp-buffer-size=4 -Wformat -Wformat-security -m32 -msse -msse2 -march=i686 -mfpmath=sse -mtune=generic -fno-omit-frame-pointer -fasynchronous-unwind-tables -std=gnu++0x -fPIC -fvisibility=hidden -fvisibility-inlines-hidden -Wall -W -D_REENTRANT -fPIC\n")

	//fmt.Fprint(bb, "#cgo CXXFLAGS: -DQT_NO_DEBUG")
	for _, m := range libs {
		fmt.Fprintf(bb, " -DQT_%v_LIB", strings.ToUpper(m))
	}
	fmt.Fprint(bb, "\n")

	fmt.Fprint(bb, "#cgo CXXFLAGS: -I/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/share/qt5/mkspecs/linux-g++ -isystem /srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/include -isystem /srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/include/sailfishapp -isystem /srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/include/mdeclarativecache5 -isystem /srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/include/qt5\n")

	fmt.Fprint(bb, "#cgo CXXFLAGS:")
	for _, m := range libs {
		fmt.Fprintf(bb, " -isystem /srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/include/qt5/Qt%v", m)
	}
	fmt.Fprint(bb, "\n\n")

	fmt.Fprint(bb, "#cgo LDFLAGS: -Wl,-O1 -Wl,-rpath-link,/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/lib -Wl,-rpath-link,/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/lib -Wl,-rpath-link,/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/lib/pulseaudio\n")

	fmt.Fprint(bb, "#cgo LDFLAGS: -rdynamic -L/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/lib -L/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/lib -L/srv/mer/targets/SailfishOS-"+utils.QT_SAILFISH_VERSION()+"-i486/usr/lib/pulseaudio -lsailfishapp -lmdeclarativecache5")
	for _, m := range libs {
		if !(m == "UiPlugin" || m == "Sailfish") {
			if parser.IsWhiteListedSailfishLib(m) {
				fmt.Fprintf(bb, " -lQt5%v", m)
			}
		}
	}

	fmt.Fprint(bb, " -lGLESv2 -lpthread\n")
	fmt.Fprint(bb, "*/\n")

	fmt.Fprint(bb, "import \"C\"\n")

	var tmp = strings.Replace(bb.String(), "${BUILDTARGET}", "sailfish_emulator", -1)

	switch {
	case mode == RCC:
		{
			utils.Save(filepath.Join(mocPath, "rcc_cgo_sailfish_emulator_linux_386.go"), tmp)
		}

	case mode == MOC:
		{
			utils.Save(filepath.Join(mocPath, "moc_cgo_sailfish_emulator_linux_386.go"), tmp)
		}

	case mode == MINIMAL:
		{
			utils.Save(utils.GoQtPkgPath(strings.ToLower(module), "minimal_cgo_sailfish_emulator_linux_386.go"), tmp)
		}

	default:
		{
			utils.Save(utils.GoQtPkgPath(strings.ToLower(module), "cgo_sailfish_emulator_linux_386.go"), tmp)
		}
	}

	tmp = strings.Replace(bb.String(), "${BUILDTARGET}", "sailfish", -1)
	tmp = strings.Replace(tmp, "-m32 -msse -msse2 -march=i686 -mfpmath=sse -mtune=generic -fno-omit-frame-pointer -fasynchronous-unwind-tables", "-fmessage-length=0 -march=armv7-a -mfloat-abi=hard -mfpu=neon -mthumb -Wno-psabi", -1)
	tmp = strings.Replace(tmp, "i486", "armv7hl", -1)

	switch {
	case mode == RCC:
		{
			utils.Save(filepath.Join(mocPath, "rcc_cgo_sailfish_linux_arm.go"), tmp)
		}

	case mode == MOC:
		{
			utils.Save(filepath.Join(mocPath, "moc_cgo_sailfish_linux_arm.go"), tmp)
		}

	case mode == MINIMAL:
		{
			utils.Save(utils.GoQtPkgPath(strings.ToLower(module), "minimal_cgo_sailfish_linux_arm.go"), tmp)
		}

	default:
		{
			utils.Save(utils.GoQtPkgPath(strings.ToLower(module), "cgo_sailfish_linux_arm.go"), tmp)
		}
	}
}

func cgoAsteroid(module, mocPath string, mode int, pkg string) {
	var (
		bb   = new(bytes.Buffer)
		libs = cleanLibs(module, mode)
	)
	defer bb.Reset()

	fmt.Fprintf(bb, "// +build ${BUILDTARGET}%v\n\n", func() string {
		if mode == MINIMAL {
			return ",minimal"
		}
		if mode == MOC {
			return ""
		}
		return ",!minimal"
	}())

	fmt.Fprintf(bb, "package %v\n\n", func() string {
		if mode == MOC {
			return pkg
		}
		return strings.ToLower(module)
	}())
	fmt.Fprint(bb, "/*\n")

	fmt.Fprint(bb, "#cgo CFLAGS: -pipe -O2 -g -pipe -Wall -Wp,-D_FORTIFY_SOURCE=2 -feliminate-unused-debug-types -fexceptions -fstack-protector --param=ssp-buffer-size=4 -Wformat -Wformat-security -fmessage-length=0 -march=armv7ve -mfloat-abi=softfp -mfpu=neon -mthumb -Wno-psabi -fPIC -fvisibility=hidden -Wall -W -D_REENTRANT -fPIE\n")
	fmt.Fprint(bb, "#cgo CXXFLAGS: -pipe -O2 -g -pipe -Wall -Wp,-D_FORTIFY_SOURCE=2 -feliminate-unused-debug-types -fexceptions -fstack-protector --param=ssp-buffer-size=4 -Wformat -Wformat-security -fmessage-length=0 -march=armv7ve -mfloat-abi=softfp -mfpu=neon -mthumb -Wno-psabi -fPIC -fvisibility=hidden -Wall -W -D_REENTRANT -fPIE\n")

	//fmt.Fprint(bb, "#cgo CXXFLAGS: -DQT_NO_DEBUG")
	for _, m := range libs {
		fmt.Fprintf(bb, " -DQT_%v_LIB", strings.ToUpper(m))
	}
	fmt.Fprint(bb, "\n")

	fmt.Fprintf(bb, "#cgo CXXFLAGS: -I%[1]s/usr/include/c++/6.2.0/arm-oe-linux-gnueabi -I%[1]s/usr/include/c++/6.2.0  -I%[1]s/usr/lib/mkspecs -I%[1]s/usr/include -I%[1]s/usr/include/mdeclarativecache5 -I%[1]s/usr/include/resource/qt5\n", os.Getenv("OECORE_TARGET_SYSROOT"))

	fmt.Fprint(bb, "#cgo CXXFLAGS:")
	for _, m := range libs {
		fmt.Fprintf(bb, " -I%s/usr/include/Qt%v", os.Getenv("OECORE_TARGET_SYSROOT"), m)
	}
	fmt.Fprint(bb, "\n\n")

	fmt.Fprintf(bb, "#cgo LDFLAGS: -Wl,-O1 -Wl,-rpath-link,%[1]s/usr/lib -Wl,-rpath-link,%[1]s/lib\n", os.Getenv("OECORE_TARGET_SYSROOT"))

	fmt.Fprintf(bb, "#cgo LDFLAGS: -rdynamic -L%[1]s/usr/lib -L%[1]s/lib -lmdeclarativecache5", os.Getenv("OECORE_TARGET_SYSROOT"))
	for _, m := range libs {
		if m != "UiPlugin" {
			if parser.IsWhiteListedSailfishLib(m) {
				fmt.Fprintf(bb, " -lQt5%v", m)
			}
		}
	}

	fmt.Fprint(bb, " -lGLESv2 -lpthread\n")
	fmt.Fprint(bb, "*/\n")

	fmt.Fprint(bb, "import \"C\"\n")

	var tmp = strings.Replace(bb.String(), "${BUILDTARGET}", "asteroid", -1)
	tmp = strings.Replace(tmp, "i486", "armv7ve", -1)

	switch {
	case mode == MOC:
		{
			utils.Save(filepath.Join(mocPath, "moc_cgo_asteroid_linux_arm.go"), tmp)
		}

	case mode == MINIMAL:
		{
			utils.Save(utils.GoQtPkgPath(strings.ToLower(module), "minimal_cgo_asteroid_linux_arm.go"), tmp)
		}

	default:
		{
			utils.Save(utils.GoQtPkgPath(strings.ToLower(module), "cgo_asteroid_linux_arm.go"), tmp)
		}
	}
}
