package moc

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/tools/imports"

	"github.com/peterq/pan-light/qt/tool-chain/binding/parser"
	"github.com/peterq/pan-light/qt/tool-chain/binding/templater"

	"github.com/peterq/pan-light/qt/tool-chain/cmd"
	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

var (
	done   = make(map[string]int)
	rootWg = new(sync.WaitGroup)

	goFilesCache      = make(map[string][]string)
	goFilesCacheMutex = new(sync.Mutex)

	goImportsCache      = make(map[string][]string)
	goImportsCacheMutex = new(sync.Mutex)

	goMocImportsCache      = make(map[string][]string)
	goMocImportsCacheMutex = new(sync.Mutex)

	gL int

	ResourceNames      = make(map[string]string)
	ResourceNamesMutex = new(sync.Mutex)

	goNameCache      = make(map[string]string)
	goNameCacheMutex = new(sync.Mutex)

	goDirCache      = make(map[string]string)
	goDirCacheMutex = new(sync.Mutex)
)

func Moc(path, target, tags string, fast, slow bool) {
	parser.State.Target = target
	if utils.UseGOMOD(path) {
		if !utils.ExistsDir(filepath.Join(path, "vendor")) {
			cmd := exec.Command("go", "mod", "vendor")
			cmd.Dir = path
			utils.RunCmd(cmd, "go mod vendor")
		}
	}

	moc(path, target, tags, fast, slow, true, -1, false)
}

func moc(path, target, tags string, fast, slow, root bool, l int, dirty bool) {
	utils.Log.WithField("path", path).WithField("target", target).Debug("start Moc")

	if !dirty {
		env, tagsEnv, _, _ := cmd.BuildEnv(target, "", "")
		scmd := utils.GoList("'{{.Stale}}':'{{.StaleReason}}'")
		scmd.Dir = path

		if !fast && !utils.QT_FAT() {
			tagsEnv = append(tagsEnv, "minimal")
		}
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
			utils.Log.WithField("path", path).Debug("skipping already cached moc")
			return
		}
		dirty = true
	}

	l++
	if l <= gL {
		//TODO: clear parser.State.GoClassMap ?
		parser.State.MocImports = make(map[string]struct{})
	}
	gL = l
	//fmt.Println(l, strings.Repeat(" ", l)+strings.Replace(path, utils.MustGoPath()+"/src/", "", -1))

	if root {
		ngr := 50
		if slow {
			ngr = 1
		}
		defer rootWg.Wait()
		wg := new(sync.WaitGroup)
		wc := make(chan bool, ngr)
		allImports := append([]string{path}, cmd.GetImports(path, target, tags, 0, false, true)...)

		wg.Add(len(allImports) * 2)
		for _, path := range allImports {
			wc <- true
			go func(path string) {
				files := cmd.GetGoFiles(path, target, tags)
				goFilesCacheMutex.Lock()
				goFilesCache[path] = files
				goFilesCacheMutex.Unlock()
				<-wc
				wg.Done()
			}(path)

			wc <- true
			go func(path string) {
				imports := cmd.GetImports(path, target, tags, 0, true, true)
				goImportsCacheMutex.Lock()
				goImportsCache[path] = imports
				goImportsCacheMutex.Unlock()
				<-wc
				wg.Done()
			}(path)
		}
		wg.Wait()
	}

	var (
		classes      []*parser.Class
		otherclasses []*parser.Class
		pkg          string
	)

	for i, spath := range append([]string{path}, goImportsCache[path]...) {
		if _, ok := done[spath]; !ok && i > 0 && !fast {
			done[spath] = i
			moc(spath, target, tags, false, slow, false, l, dirty)
		}

		for _, fpath := range goFilesCache[spath] {
			cls, ipkg, err := parse(fpath)
			if pkg == "" || strings.TrimSuffix(pkg, "_test") == ipkg {
				pkg = ipkg
			}
			if err != nil {
				utils.Log.WithError(err)
				continue
			}
			if cls == nil {
				continue
			}

			if spath == path {
				classes = append(classes, cls...)
				continue
			}

			importAs := "custom_" + ipkg + "_" + cls[0].Hash() + "m"
			if strings.Contains(spath, "/vendor/") {
				importAs = ipkg
			}

			(&parser.Module{
				Project:   importAs,
				Pkg:       spath,
				Namespace: &parser.Namespace{Classes: cls},
			}).Prepare()

			otherclasses = append(otherclasses, cls...)
		}
	}

	utils.Log.WithField("path", path).Debugln("found", len(classes), "moc structs")
	if len(classes) == 0 {
		return
	}

	if _, ok := parser.State.ClassMap["QObject"]; !ok {
		parser.LoadModules()
	} else {
		utils.Log.Debug("modules already cached")
	}

	//find valid base classes
	for _, c := range classes {
		for _, bc := range c.GetBases() {
			var has bool
			for _, c := range append(classes, otherclasses...) {
				if c.Name == bc {
					has = true
					break
				}
			}
			if !has {
				for _, c := range parser.State.ClassMap {
					if c.Name != bc {
						has = true
						break
					}
				}
			}
			if has {
				c.Bases = bc
				break
			}
		}
	}

	//TODO: tool-chain binding functions rely on the prepared state
	m := &parser.Module{
		Project:   parser.MOC,
		Namespace: &parser.Namespace{Classes: classes},
	}

	//copy properties + signals + slots
	utils.Log.Debug("start copy properties + signals + slots")
	for range append(m.Namespace.Classes, otherclasses...) {
		for _, c := range append(m.Namespace.Classes, otherclasses...) {
			bc, ok := parser.State.ClassMap[c.Bases]
			if !ok || bc.Pkg == "" {
				continue
			}

			for _, bcp := range bc.Properties {
				var has bool
				for _, cp := range c.Properties {
					if cp.Name == bcp.Name {
						has = true
						break
					}
				}
				if !has {
					bcp.IsMocSynthetic = true
					c.Properties = append(c.Properties, bcp)
				}
			}

			for _, bcf := range bc.Functions {
				if !(bcf.Meta == parser.SIGNAL || bcf.Meta == parser.SLOT) {
					continue
				}

				f := *bcf
				f.Name = strings.Replace(f.Name, bcf.ClassName(), c.Name, -1)
				f.Fullname = strings.Replace(f.Fullname, bcf.ClassName(), c.Name, -1)
				if !c.HasFunctionWithNameAndOverloadNumber(f.Name, f.OverloadNumber) {
					c.Functions = append(c.Functions, &f)
				}
			}
		}
	}
	utils.Log.Debug("done copy properties + signals + slots")

	m.Prepare()

	var remaining int
	for _, class := range m.Namespace.Classes {
		if _, failed := class.GetAllBasesRecursiveCheckFailed(0); failed || !class.IsSubClassOfQObject() {
			delete(parser.State.ClassMap, class.Name)
		} else {
			remaining++
		}
	}
	utils.Log.WithField("path", path).Debugln("found", remaining, "remaining moc structs")
	if remaining == 0 {
		return
	}

	utils.Log.Debug("start converting types")
	for _, c := range m.Namespace.Classes {
		for _, f := range c.Functions {
			if f.NoMocDeduce {
				continue
			}
			for _, p := range f.Parameters {
				p.Value, p.PureGoType = cppTypeFromGoType(f, p.Value, c)
			}
			f.Output, f.PureGoOutput = cppTypeFromGoType(f, f.Output, c)
		}
		utils.Log.Debug("done converting func types")
		for _, p := range c.Properties {
			p.Output, p.PureGoType = cppTypeFromGoType(nil, p.Output, c)
		}
		utils.Log.Debug("done converting property types")

		//TODO: needed because only now the values are decuded to c++ types
		c.FixGenericHelper()
	}
	utils.Log.Debug("done converting types")

	//copy constructor and destructor
	utils.Log.Debug("start copy structors")
	for !hasStructors(m) {
		for _, c := range append(m.Namespace.Classes, otherclasses...) {
			bc, ok := parser.State.ClassMap[c.Bases]
			if !ok {
				continue
			}
			for _, bcf := range bc.Functions {
				if !(bcf.Meta == parser.CONSTRUCTOR || bcf.Meta == parser.DESTRUCTOR) {
					continue
				}
				f := *bcf
				f.Name = strings.Replace(f.Name, bcf.ClassName(), c.Name, -1)
				f.Fullname = strings.Replace(f.Fullname, bcf.ClassName(), c.Name, -1)
				if !c.HasFunctionWithNameAndOverloadNumber(f.Name, f.OverloadNumber) {
					c.Functions = append(c.Functions, &f)
				}
			}
		}
	}
	utils.Log.Debug("done copy structors")

	goFilesCacheMutex.Lock()
	files := goFilesCache[path]
	goFilesCacheMutex.Unlock()
	parseNonMocDeps(files)

	ResourceNamesMutex.Lock()
	ResourceNames[filepath.Join(path, "moc.cpp")] = ""
	ResourceNamesMutex.Unlock()

	env, tagsEnv, _, _ := cmd.BuildEnv(target, "", "")
	scmd := utils.GoList("'{{.Stale}}':'{{.StaleReason}}'")
	scmd.Dir = path

	if !fast && !utils.QT_FAT() {
		tagsEnv = append(tagsEnv, "minimal")
	}
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
		utils.Log.WithField("path", path).Debug("skipping already cached moc")
	} else {
		if err := utils.SaveBytes(filepath.Join(path, "moc.cpp"), templater.CppTemplate(parser.MOC, templater.MOC, target, tags)); err != nil {
			return
		}

		if err := utils.SaveBytes(filepath.Join(path, "moc.h"), templater.HTemplate(parser.MOC, templater.MOC, tags)); err != nil {
			return
		}
		fixR := templater.GoTemplate(parser.MOC, false, templater.MOC, pkg, target, tags)

		rootWg.Add(1)
		go func() {
			defer rootWg.Done()

			var fix []byte
			var err error
			for i := 0; i < 5; i++ {
				fix, err = imports.Process("moc.go", fixR, nil)
				if err != nil {
					log.Println("here")
					utils.Log.WithError(err).Error("failed to fix go imports")
					fix = fixR
				}
			}

			if err := utils.SaveBytes(filepath.Join(path, "moc.go"), fix); err != nil {
				return
			}
		}()

		var libs []string
		parser.LibDepsMutex.Lock()
		libs = parser.LibDeps[parser.MOC]
		parser.LibDepsMutex.Unlock()

		rootWg.Add(1)
		go func() {
			templater.CgoTemplateSafe(parser.MOC, path, target, templater.MOC, pkg, tags, libs)
			rootWg.Done()
		}()

		rootWg.Add(1)
		go func() {
			utils.RunCmd(exec.Command(utils.ToolPath("moc", target), filepath.Join(path, "moc.cpp"), "-o", filepath.Join(path, "moc_moc.h")), "run moc")
			if tags != "" {
				utils.Save(filepath.Join(path, "moc_moc.h"), "// +build "+tags+"\n\n"+utils.Load(filepath.Join(path, "moc_moc.h")))
			}

			if utils.QT_DOCKER() {
				if idug, ok := os.LookupEnv("IDUG"); ok {
					utils.RunCmd(exec.Command("chown", "-R", idug, path), "chown files to user")
				}
			}
			rootWg.Done()
		}()
	}

	//TODO: cleanup state -->
	for _, c := range parser.State.ClassMap {
		if c.Module == parser.MOC {
			delete(parser.State.ClassMap, c.Name)
		}
	}
	parser.LibDepsMutex.Lock()
	parser.LibDeps[parser.MOC] = make([]string, 0)
	parser.LibDepsMutex.Unlock()
	//<--
}

func parse(path string) ([]*parser.Class, string, error) {
	if base := filepath.Base(path); strings.HasPrefix(base, "rcc") || strings.HasPrefix(base, "moc") {
		return nil, "", errors.New("file contains no moc structs")
	}

	if strings.HasPrefix(path, filepath.Join(runtime.GOROOT(), "src")) {
		return nil, "", errors.New("path is in GOROOT/src")
	}

	utils.Log.WithField("path", path).Debug("parse")

	src, err := ioutil.ReadFile(path)
	if err != nil {
		utils.Log.WithField("path", path).WithError(err).Debug("failed to read file")
		return nil, "", err
	}

	file, err := goparser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		utils.Log.WithField("path", path).WithError(err).Debug("failed to parse file")
		return nil, "", err
	}

	//TODO: cache moc struct hashes in moc.go files to indentify staled moc structs in staled go packages

	var classes []*parser.Class
	for _, d := range file.Decls {
		typeDecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, s := range typeDecl.Specs {
			typeSpec, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structDecl, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			class := &parser.Class{
				Access: "public",
				Module: parser.MOC,
				Name:   typeSpec.Name.String(),
				Status: "public",
				Path:   filepath.Dir(path),
			}

			//collect possible base classes
			var bases []string
			for _, field := range structDecl.Fields.List {
				value := string(src[field.Pos()-1 : field.End()-1])
				if value != "" && !strings.Contains(value, " ") && !strings.HasPrefix(value, "*") {
					if strings.Contains(value, ".") {
						value = strings.Split(value, ".")[1]
					}
					bases = append(bases, value)
				}

				if len(field.Names) > 0 {
					if field.Tag == nil {
						continue
					}
					tag := strings.Replace(strings.Replace(field.Tag.Value, "\"", "", -1), "`", "", -1)

					var meta string
					switch {
					case strings.HasPrefix(tag, "signal:"):
						meta = parser.SIGNAL
					case strings.HasPrefix(tag, "slot:"):
						meta = parser.SLOT
					case strings.HasPrefix(tag, "property:"):
						meta = parser.PROP
					case strings.HasPrefix(tag, "constructor:"): //TODO: more advanced constructor support (multiple constructors, custom inputs, error output, custom naming, ...)
						meta = parser.CONSTRUCTOR
					default:
						continue
					}

					//TODO: sync, async, lazy, ...
					//TODO: whole class shims
					//TODO: multi targets
					//TODO: private
					//TODO: qml register tag
					var (
						auto    int
						target  string
						inbound bool
					)
					if (meta == parser.SIGNAL || meta == parser.SLOT || meta == parser.PROP) && (strings.Contains(tag, ",->") || strings.Contains(tag, ",auto") || strings.Contains(tag, ",<-")) {

						autoTag := strings.Split(tag, ",")[1]

						if strings.Contains(tag, ",->") || strings.Contains(tag, ",auto") {
							auto = 1
							autoTag = strings.TrimPrefix(autoTag, "->")
							autoTag = strings.TrimPrefix(autoTag, "auto")
						} else {
							auto = 2
							autoTag = strings.TrimPrefix(autoTag, "<-")
						}

						if strings.Contains(autoTag, "(") {
							if !strings.HasPrefix(autoTag, "(this.") {
								//TODO: concurrent?
								var found bool
								for _, imp := range file.Imports {
									if strings.Contains(autoTag, "("+imp.Name.String()+".") {
										var name string
										goNameCacheMutex.Lock()
										if n, ok := goNameCache[imp.Path.Value]; ok {
											name = n
										} else {
											name = strings.TrimSpace(utils.RunCmd(utils.GoList("{{.Name}}", strings.Replace(imp.Path.Value, "\"", "", -1)), "get import name"))
											goNameCache[imp.Path.Value] = name
										}
										goNameCacheMutex.Unlock()

										var dir string
										goDirCacheMutex.Lock()
										if d, ok := goDirCache[imp.Path.Value]; ok {
											dir = d
										} else {
											dir = strings.TrimSpace(utils.RunCmd(utils.GoList("{{.Dir}}", strings.Replace(imp.Path.Value, "\"", "", -1)), "get import dir"))
											goDirCache[imp.Path.Value] = dir
										}
										goDirCacheMutex.Unlock()

										h := sha1.New()
										h.Write([]byte(dir))

										mname := "custom_" + name + "_" + hex.EncodeToString(h.Sum(nil)[:3]) + "m"
										goMocImportsCacheMutex.Lock()
										goMocImportsCache[path] = append(goMocImportsCache[path], fmt.Sprintf("%v %v", mname, imp.Path.Value))
										goMocImportsCacheMutex.Unlock()

										autoTag = strings.Replace(autoTag, "("+imp.Name.String()+".", "("+mname+".", -1)
										found = true
										break
									}
								}

								if !found {
									for _, imp := range file.Imports {
										if imp.Name.String() != "<nil>" && imp.Name.String() != "_" {
											continue
										}

										var name string
										goNameCacheMutex.Lock()
										if n, ok := goNameCache[imp.Path.Value]; ok {
											name = n
										} else {
											name = strings.TrimSpace(utils.RunCmd(utils.GoList("{{.Name}}", strings.Replace(imp.Path.Value, "\"", "", -1)), "get import name"))
											goNameCache[imp.Path.Value] = name
										}
										goNameCacheMutex.Unlock()

										if !strings.Contains(autoTag, "("+name+".") {
											continue
										}

										var dir string
										goDirCacheMutex.Lock()
										if d, ok := goDirCache[imp.Path.Value]; ok {
											dir = d
										} else {
											dir = strings.TrimSpace(utils.RunCmd(utils.GoList("{{.Dir}}", strings.Replace(imp.Path.Value, "\"", "", -1)), "get import dir"))
											goDirCache[imp.Path.Value] = dir
										}
										goDirCacheMutex.Unlock()

										h := sha1.New()
										h.Write([]byte(dir))

										mname := "custom_" + name + "_" + hex.EncodeToString(h.Sum(nil)[:3]) + "m"
										goMocImportsCacheMutex.Lock()
										goMocImportsCache[path] = append(goMocImportsCache[path], fmt.Sprintf("%v %v", mname, imp.Path.Value))
										goMocImportsCacheMutex.Unlock()

										autoTag = strings.Replace(autoTag, "("+name+".", "("+mname+".", -1)
										break
									}
								}
							}
							target = strings.TrimSuffix(strings.TrimPrefix(autoTag, "("), ")")
						}

						tag = strings.Split(tag, ",")[0]
					}

					name := strings.Split(tag, ":")[1]
					if name[:1] != strings.ToLower(name[:1]) {
						inbound = true
						name = strings.ToLower(name[:1]) + strings.TrimPrefix(name[1:], name[:1])
					}

					typ := string(src[field.Type.Pos()-1 : field.Type.End()-1])
					switch meta {
					case parser.SIGNAL, parser.SLOT:
						f := &parser.Function{
							Access:        "public",
							Fullname:      fmt.Sprintf("%v::%v", class.Name, name),
							Meta:          meta,
							Name:          name,
							Status:        "public",
							Virtual:       parser.PURE,
							Signature:     "()",
							Output:        "void",
							Parameters:    parameters(typ),
							IsMocFunction: true,
							Connect:       auto,
							Target:        target,
							Inbound:       inbound,
						}
						if meta == parser.SLOT {
							if strings.Contains(typ, ") ") {
								f.Output = strings.Split(typ, ") ")[1]
							} else if strings.Contains(typ, ")") {
								f.Output = strings.Split(typ, ")")[1]
							}
							if strings.HasPrefix(f.Output, "(") {
								f.Output = strings.TrimPrefix(f.Output, "(")
								f.Output = strings.TrimSuffix(f.Output, ")")
								f.Output = strings.Join(strings.Split(f.Output, " ")[1:], " ")
							}
						}
						if len(f.Parameters[0].Value) == 0 {
							f.Parameters = make([]*parser.Parameter, 0)
						}
						class.Functions = append(class.Functions, f)

					case parser.PROP:
						class.Properties = append(class.Properties,
							&parser.Variable{
								Access:         "public",
								Fullname:       fmt.Sprintf("%v::%v", class.Name, strings.Split(tag, ":")[1]),
								Name:           strings.Split(tag, ":")[1],
								Status:         "public",
								Output:         typ,
								Connect:        auto,
								ConnectGet:     strings.Contains(field.Tag.Value, ",get"),
								ConnectSet:     strings.Contains(field.Tag.Value, ",set"),
								ConnectChanged: strings.Contains(field.Tag.Value, ",changed"),
								Target:         target,
								Inbound:        inbound,
							})

					case parser.CONSTRUCTOR:
						class.Constructors = append(class.Constructors, strings.Split(tag, ":")[1])
					}

				}
			}
			class.Bases = strings.Join(bases, ",")
			if len(class.Properties) != 0 || len(class.Functions) != 0 ||
				len(class.Constructors) != 0 || len(class.GetBases()) != 0 {
				classes = append(classes, class)
			} else {
				ipkg := file.Name.String()
				importAs := "custom_" + ipkg + "_" + class.Hash() + "m"
				if strings.Contains(path, "/vendor/") {
					importAs = ipkg
				}
				class.Module = importAs
				class.Pkg = filepath.Dir(path)
				parser.State.GoClassMap[class.Name] = class
			}
		}
	}

	return classes, file.Name.String(), nil
}

func parameters(tag string) []*parser.Parameter {
	if !strings.Contains(tag, "(") {
		return nil
	}

	tag = strings.TrimPrefix(tag, "func(")

	switch {
	case strings.Contains(tag, ") ("):
		tag = strings.Split(tag, ") (")[0]

	case strings.Contains(tag, ") func"):
		tag = strings.Split(tag, ") func")[0]

	case strings.Contains(tag, ") "):
		tag = strings.Split(tag, ") ")[0]

	default:
		tag = strings.TrimSuffix(tag, ")")
	}

	tag = strings.Replace(tag, " func", ";func", -1)
	tag = strings.Replace(tag, ";", " ", 1)
	tag = strings.Replace(tag, "<-chan ", "$RC_", -1)
	tag = strings.Replace(tag, "chan<- ", "$WC_", -1)
	tag = strings.Replace(tag, "chan ", "$C_", -1)

	var lv string
	var o []*parser.Parameter
	pairs := strings.Split(tag, ",")
	for i := len(pairs) - 1; i >= 0; i-- {

		var pOut *parser.Parameter
		ps := strings.Split(strings.TrimSpace(pairs[i]), " ")
		if len(ps) == 1 {
			pOut = &parser.Parameter{Name: fmt.Sprintf("v%v", i), Value: ps[0]}
			if lv != "" {
				pOut.Name = pOut.Value
				pOut.Value = lv
			}
		} else {
			pOut = &parser.Parameter{Name: ps[0], Value: ps[1]}
			lv = pOut.Value
		}

		o = append(o, pOut)
	}

	var ro []*parser.Parameter
	for i := len(o) - 1; i >= 0; i-- {
		o[i].Name = strings.Replace(o[i].Name, ";", " ", -1)
		o[i].Value = strings.Replace(o[i].Value, ";", " ", -1)
		ro = append(ro, o[i])
	}

	return ro
}

func cppTypeFromGoType(f *parser.Function, t string, class *parser.Class) (string, string) {
	tOld := t //TODO: also for differentiation of QVariant and *QVariant
	//t = strings.TrimPrefix(t, "*")

	if strings.Count(t, "[") == 1 || strings.HasSuffix(t, "][]string") {
		//TODO: multidimensional array and nested maps
		if strings.HasPrefix(t, "[]") && t != "[]string" {
			o, pureGoType := cppTypeFromGoType(f, strings.TrimPrefix(t, "[]"), class)
			if pureGoType == "" {
				return fmt.Sprintf("QList<%v>", o), ""
			} else if strings.Contains(pureGoType, "error") {
				return fmt.Sprintf("QList<%v>", o), t
			}
		}
		if strings.HasPrefix(t, "map[") {
			var head = fmt.Sprintf("map[%v]", strings.Split(strings.TrimPrefix(t, "map["), "]")[0])
			o1, pureGoType1 := cppTypeFromGoType(f, strings.Split(strings.TrimPrefix(t, "map["), "]")[0], class)
			o2, pureGoType2 := cppTypeFromGoType(f, strings.TrimPrefix(t, head), class)
			if pureGoType1 == "" && pureGoType2 == "" {
				return fmt.Sprintf("QMap<%v, %v>", o1, o2), ""
			} else if strings.Contains(pureGoType1, "error") || strings.Contains(pureGoType2, "error") {
				return fmt.Sprintf("QMap<%v, %v>", o1, o2), t
			}
		}
	}

	switch t {
	case "string":
		return "QString", ""
	case "error":
		return "QString", t
	case "[]string":
		return "QStringList", ""
	case "bool":
		return "bool", ""
	case "int8":
		return "qint8", ""
	case "uint8":
		return "quint8", ""
	case "int16":
		return "qint16", ""
	case "uint16":
		return "quint16", ""
	case "int", "int32":
		return "qint32", ""
	case "uint", "uint32":
		return "quint32", ""
	case "int64":
		return "qint64", ""
	case "uint64":
		return "quint64", ""
	case "float32":
		return "float", ""
	case "float64":
		return "qreal", ""
	case "uintptr":
		return "quintptr", ""
	case "unsafe.Pointer":
		return "void*", ""
	}

	if strings.Contains(t, ".") {
		t = strings.Split(t, ".")[1]
	}

	if strings.Contains(t, "__") {
		return strings.Replace(t, "_", ":", -1), ""
	}

	//TODO: differentiate between QVariant and *QVariant
	ttmp := strings.TrimPrefix(t, "*")
	if _, ok := parser.State.ClassMap[ttmp]; ok {
		if parser.State.ClassMap[ttmp].IsSubClassOfQObject() /*TODO: || f == nil && strings.HasPrefix(tOLD, "*")*/ {
			return ttmp + "*", ""
		}
		return ttmp, ""
	}

	//TODO: *_ITF support

	if tOld == "" || tOld == "void" {
		return "void", ""
	}

	tOld = strings.Replace(tOld, "$RC_", "<-chan ", -1)
	tOld = strings.Replace(tOld, "$WC_", "chan<- ", -1)
	tOld = strings.Replace(tOld, "$C_", "chan ", -1)

	//TODO: directly resolve moc pkgs imports in parse ?
	ttOld := strings.TrimPrefix(tOld, "*")
	if strings.Contains(tOld, ".") {
		ttOld = strings.Split(ttOld, ".")[1]
	}
	if c, ok := parser.State.GoClassMap[ttOld]; ok {
		if c.Path != class.Path {
			pos := c.Module + "." + ttOld
			if strings.HasPrefix(tOld, "*") {
				pos = "*" + pos
			}
			return "quintptr", pos
		}
	}

	return "quintptr", tOld
}

func hasStructors(m *parser.Module) bool {
	for _, c := range m.Namespace.Classes {
		if !c.IsSubClassOfQObject() {
			continue
		}
		if !c.HasConstructor() /*|| !c.HasDestructor()*/ {
			return false
		}
	}
	return true
}

//TODO: replace MocImports hack with GoClassMap, needs parse() to properly detect all pure Go types, structs, interfaces, ...
func parseNonMocDeps(files []string) {
	utils.Log.Debug("parseNonMocDeps")

	wg := new(sync.WaitGroup)
	wc := make(chan bool, 50)

	for _, fpath := range files {
		if base := filepath.Base(fpath); strings.HasPrefix(base, "rcc") || strings.HasPrefix(base, "moc") {
			continue
		}

		wg.Add(1)
		wc <- true
		go func(fpath string) {

			goMocImportsCacheMutex.Lock()
			imps, ok := goMocImportsCache[fpath]
			goMocImportsCacheMutex.Unlock()

			if !ok {
				utils.Log.Debugf("parse non moc deps: %v", fpath)
				file, err := goparser.ParseFile(token.NewFileSet(), fpath, nil, 0)
				if err != nil {
					<-wc
					wg.Done()
					return
				}

				for _, i := range file.Imports {
					if i.Path.Value == "\"C\"" {
						continue
					}

					if cmd.IsStdPkg(strings.Replace(i.Path.Value, "\"", "", -1)) {
						if i.Name != nil {
							goMocImportsCacheMutex.Lock()
							goMocImportsCache[fpath] = append(goMocImportsCache[fpath], fmt.Sprintf("%v %v", i.Name.String(), i.Path.Value))
							goMocImportsCacheMutex.Unlock()
						} else {
							goMocImportsCacheMutex.Lock()
							goMocImportsCache[fpath] = append(goMocImportsCache[fpath], i.Path.Value)
							goMocImportsCacheMutex.Unlock()
						}
					}
				}

				goMocImportsCacheMutex.Lock()
				imps = goMocImportsCache[fpath]
				goMocImportsCacheMutex.Unlock()
			}

			for _, path := range imps {
				goMocImportsCacheMutex.Lock()
				parser.State.MocImports[path] = struct{}{}
				goMocImportsCacheMutex.Unlock()
			}

			<-wc
			wg.Done()
		}(fpath)
	}

	wg.Wait()
}
