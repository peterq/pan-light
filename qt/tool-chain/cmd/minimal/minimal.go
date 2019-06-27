package minimal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/peterq/pan-light/qt/tool-chain/binding/converter"
	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
	"github.com/peterq/pan-light/qt/tool-chain/binding/templater"

	"github.com/peterq/pan-light/qt/tool-chain/cmd"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func Minimal(path, target, tags string) {
	parser.State.Target = target
	if utils.UseGOMOD(path) {
		if !utils.ExistsDir(filepath.Join(path, "vendor")) {
			cmd := exec.Command("go", "mod", "vendor")
			cmd.Dir = path
			utils.RunCmd(cmd, "go mod vendor")
		}
	}

	env, tagsEnv, _, _ := cmd.BuildEnv(target, "", "")
	scmd := utils.GoList("'{{.Stale}}':'{{.StaleReason}}'")
	scmd.Dir = path

	tagsEnv = append(tagsEnv, "minimal")

	if tags != "" {
		tagsEnv = append(tagsEnv, strings.Split(tags, " ")...)
	}
	scmd.Args = append(scmd.Args, fmt.Sprintf("-tags=\"%v\"", strings.Join(tagsEnv, "\" \"")))

	if target != runtime.GOOS {
		scmd.Args = append(scmd.Args, []string{"-pkgdir", filepath.Join(utils.MustGoPath(), "pkg", fmt.Sprintf("%v_%v_%v", strings.Replace(target, "-", "_", -1), env["GOOS"], env["GOARCH"]))}...)
	}

	for key, value := range env {
		scmd.Env = append(scmd.Env, fmt.Sprintf("%v=%v", key, value))
	}

	if out := utils.RunCmdOptional(scmd, fmt.Sprintf("go check stale for %v on %v", target, runtime.GOOS)); strings.Contains(out, "but available in build cache") || strings.Contains(out, "false") {
		utils.Log.WithField("path", path).Debug("skipping already cached minimal")
		return
	}

	utils.Log.WithField("path", path).WithField("target", target).Debug("start Minimal")

	//TODO: cleanup state from moc for minimal first -->
	for _, c := range parser.State.ClassMap {
		if c.Module == parser.MOC || strings.HasPrefix(c.Module, "custom_") {
			delete(parser.State.ClassMap, c.Name)
		}
	}
	parser.LibDeps[parser.MOC] = make([]string, 0)
	if target == "js" || target == "wasm" { //TODO: REVIEW
		if parser.LibDeps["build_static"][0] == "Qml" {
			parser.LibDeps["build_static"] = parser.LibDeps["build_static"][1:]
		}
	} else {
		parser.LibDeps["build_static"] = []string{"Qml"}
	}
	//<--

	wg := new(sync.WaitGroup)
	wc := make(chan bool, 50)

	var files []string
	fileMutex := new(sync.Mutex)

	allImports := append([]string{path}, cmd.GetImports(path, target, tags, 0, false, false)...)
	wg.Add(len(allImports))
	for _, path := range allImports {
		wc <- true
		go func(path string) {
			for _, path := range cmd.GetGoFiles(path, target, tags) {
				if base := filepath.Base(path); strings.HasPrefix(base, "rcc_cgo") || strings.HasPrefix(base, "moc_cgo") {
					continue
				}
				utils.Log.WithField("path", path).Debug("analyse for minimal")
				file := utils.Load(path)
				fileMutex.Lock()
				files = append(files, file)
				fileMutex.Unlock()
			}
			if target == "js" { //TODO: wasm as well
				filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
					if err != nil || info.IsDir() {
						return err
					}
					if filepath.Ext(path) == ".js" {
						utils.Log.WithField("path", path).Debug("analyse js for minimal")
						file := utils.Load(path)
						fileMutex.Lock()
						files = append(files, file)
						fileMutex.Unlock()
					}
					return nil
				})
			}
			<-wc
			wg.Done()
		}(path)
	}
	wg.Wait()

	c := len(files)
	utils.Log.Debugln("found", c, "files to analyze")
	if c == 0 {
		return
	}

	if _, ok := parser.State.ClassMap["QObject"]; !ok {
		parser.LoadModules()
	} else {
		utils.Log.Debug("modules already cached")
	}

	//TODO: merge sailfish and asteroid
	switch target {
	case "sailfish", "sailfish-emulator":
		for _, bl := range []string{"TestCase", "QQuickWidget", "QLatin1String", "QStringRef"} {
			if c, ok := parser.State.ClassMap[bl]; ok {
				c.Export = false
				delete(parser.State.ClassMap, c.Name)
			}
		}

		for _, c := range parser.State.ClassMap {
			since, err := strconv.ParseFloat(strings.TrimPrefix(c.Since, "Qt "), 64)
			if err != nil {
				continue
			}
			if since >= 5.3 || !parser.IsWhiteListedSailfishLib(strings.TrimPrefix(c.Module, "Qt")) {
				c.Export = false
				delete(parser.State.ClassMap, c.Name)
				continue
			}

			for _, f := range c.Functions {
				since, err := strconv.ParseFloat(strings.TrimPrefix(f.Since, "Qt "), 64)
				if err != nil {
					continue
				}
				if since >= 5.3 {
					f.Export = false
				}
			}
		}

	case "asteroid":
		for _, bl := range []string{"TestCase", "QQuickWidget"} {
			if c, ok := parser.State.ClassMap[bl]; ok {
				c.Export = false
				delete(parser.State.ClassMap, c.Name)
			}
		}

		for _, c := range parser.State.ClassMap {
			since, err := strconv.ParseFloat(strings.TrimPrefix(c.Since, "Qt "), 64)
			if err != nil {
				continue
			}
			if since >= 5.7 || !parser.IsWhiteListedSailfishLib(strings.TrimPrefix(c.Module, "Qt")) {
				c.Export = false
				delete(parser.State.ClassMap, c.Name)
				continue
			}

			for _, f := range c.Functions {
				since, err := strconv.ParseFloat(strings.TrimPrefix(f.Since, "Qt "), 64)
				if err != nil {
					continue
				}
				if since >= 5.7 {
					f.Export = false
				}
			}
		}

	case "ios", "ios-simulator":
		for _, bl := range []string{"QProcess", "QProcessEnvironment"} {
			if c, ok := parser.State.ClassMap[bl]; ok {
				c.Export = false
				delete(parser.State.ClassMap, bl)
			}
		}

	case "rpi1", "rpi2", "rpi3":
		if !utils.QT_RPI() {
			break
		}
		for _, bl := range []string{"TestCase", "QQuickWidget"} {
			if c, ok := parser.State.ClassMap[bl]; ok {
				c.Export = false
				delete(parser.State.ClassMap, c.Name)
			}
		}

		for _, c := range parser.State.ClassMap {
			since, err := strconv.ParseFloat(strings.TrimPrefix(c.Since, "Qt "), 64)
			if err != nil {
				continue
			}
			if since >= 5.8 || !parser.IsWhiteListedRaspberryLib(strings.TrimPrefix(c.Module, "Qt")) {
				c.Export = false
				delete(parser.State.ClassMap, c.Name)
				continue
			}

			for _, f := range c.Functions {
				since, err := strconv.ParseFloat(strings.TrimPrefix(f.Since, "Qt "), 64)
				if err != nil {
					continue
				}
				if since >= 5.8 {
					f.Export = false
				}
			}
		}
	case "js", "wasm":
		parser.State.ClassMap["QSvgWidget"].Export = true
	}

	wg.Add(len(files))
	for _, f := range files {
		go func(f string) {
			for _, c := range parser.State.ClassMap {
				if strings.Contains(f, c.Name) &&
					strings.Contains(f, fmt.Sprintf("github.com/peterq/pan-light/qt/bindings/%v", strings.ToLower(strings.TrimPrefix(c.Module, "Qt")))) {
					exportClass(c, files)
				}
			}
			wg.Done()
		}(f)
	}
	wg.Wait()

	parser.State.ClassMap["QVariant"].Export = true

	//TODO: cleanup state
	parser.State.Minimal = true
	for _, m := range parser.GetLibs() {
		if !parser.ShouldBuildForTarget(m, target) ||
			m == "AndroidExtras" || m == "Sailfish" {
			continue
		}

		utils.MkdirAll(utils.GoQtPkgPath(strings.ToLower(m)))
		utils.SaveBytes(utils.GoQtPkgPath(strings.ToLower(m), strings.ToLower(m)+"-minimal.cpp"), templater.CppTemplate(m, templater.MINIMAL, target, ""))
		utils.SaveBytes(utils.GoQtPkgPath(strings.ToLower(m), strings.ToLower(m)+"-minimal.h"), templater.HTemplate(m, templater.MINIMAL, ""))
		utils.SaveBytes(utils.GoQtPkgPath(strings.ToLower(m), strings.ToLower(m)+"-minimal.go"), templater.GoTemplate(m, false, templater.MINIMAL, m, target, ""))
	}

	for _, m := range parser.GetLibs() {
		if !parser.ShouldBuildForTarget(m, target) ||
			m == "AndroidExtras" || m == "Sailfish" {
			continue
		}

		wg.Add(1)
		go func(m string, libs []string) {
			templater.CgoTemplateSafe(m, "", target, templater.MINIMAL, m, "", libs)
			wg.Done()
		}(m, parser.LibDeps[m])
	}
	wg.Wait()

	parser.State.Minimal = false
	for _, c := range parser.State.ClassMap {
		c.Export = false
		for _, f := range c.Functions {
			f.Export = false
		}
	}
}

func exportClass(c *parser.Class, files []string) {
	c.Lock()
	exp := c.Export
	c.Unlock()
	if exp {
		return
	}
	c.Lock()
	c.Export = true
	c.Unlock()

	for _, file := range files {
		for _, f := range c.Functions {

			switch {
			case f.Virtual == parser.IMPURE, f.Virtual == parser.PURE, f.Meta == parser.SIGNAL, f.Meta == parser.SLOT:
				for _, mode := range []string{parser.CONNECT, parser.DISCONNECT, ""} {
					f.SignalMode = mode
					if strings.Contains(file, converter.GoHeaderName(f)) {
						exportFunction(f, files)
					}
				}

			default:
				if f.Static {
					if strings.Contains(file, converter.GoHeaderName(f)) {
						exportFunction(f, files)
					}
					f.Static = false
					if strings.Contains(file, converter.GoHeaderName(f)) {
						exportFunction(f, files)
					}
					f.Static = true
				} else {
					if strings.Contains(file, converter.GoHeaderName(f)) {
						exportFunction(f, files)
					}
				}
			}

			if strings.HasPrefix(f.Name, "__") || f.Meta == parser.CONSTRUCTOR ||
				f.Meta == parser.DESTRUCTOR || f.Virtual == parser.PURE {
				exportFunction(f, files)
			}

		}
	}

	for _, b := range c.GetAllBases() {
		if c, ok := parser.State.ClassMap[b]; ok {
			exportClass(c, files)
		}
	}
}

func exportFunction(f *parser.Function, files []string) {
	if f.Export {
		return
	}
	f.Export = true

	for _, p := range f.Parameters {
		if c, ok := parser.State.ClassMap[parser.CleanValue(p.Value)]; ok {
			exportClass(c, files)
		}
		if parser.IsPackedList(p.Value) {
			if c, ok := parser.State.ClassMap[parser.UnpackedList(p.Value)]; ok {
				exportClass(c, files)
			}
		}
		if parser.IsPackedMap(p.Value) {
			key, value := parser.UnpackedMap(p.Value)
			if c, ok := parser.State.ClassMap[key]; ok {
				exportClass(c, files)
			}
			if c, ok := parser.State.ClassMap[value]; ok {
				exportClass(c, files)
			}
		}
	}

	if c, ok := parser.State.ClassMap[parser.CleanValue(f.Output)]; ok {
		exportClass(c, files)
	}
	if parser.IsPackedList(f.Output) {
		if c, ok := parser.State.ClassMap[parser.UnpackedList(f.Output)]; ok {
			exportClass(c, files)
		}
	}
	if parser.IsPackedMap(f.Output) {
		key, value := parser.UnpackedMap(f.Output)
		if c, ok := parser.State.ClassMap[key]; ok {
			exportClass(c, files)
		}
		if c, ok := parser.State.ClassMap[value]; ok {
			exportClass(c, files)
		}
	}

	if c, ok := parser.State.ClassMap[parser.CleanValue(f.Output)]; ok && f.Virtual == parser.PURE {
		for _, f := range c.Functions {
			if f.Meta == parser.CONSTRUCTOR && len(f.Parameters) == 0 {
				exportFunction(f, files)
			}
		}
	}
}
