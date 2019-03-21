package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

var (
	std      []string
	stdMutex = new(sync.Mutex)

	imported      = make(map[string]string)
	importedMutex = new(sync.Mutex)
)

func IsStdPkg(pkg string) bool {

	stdMutex.Lock()
	if std == nil {
		std = append(strings.Split(strings.TrimSpace(utils.RunCmd(exec.Command("go", "list", "std"), "go list std")), "\n"), "C")
	}
	stdMutex.Unlock()

	for _, spkg := range std {
		if pkg == spkg {
			return true
		}
	}
	return false
}

func GetImports(path, target, tagsCustom string, level int, onlyDirect, moc bool) []string {
	utils.Log.WithField("path", path).WithField("level", level).Debug("get imports")

	env, tags, _, _ := BuildEnv(target, "", "")

	stdMutex.Lock()
	if std == nil {
		std = append(strings.Split(strings.TrimSpace(utils.RunCmd(exec.Command("go", "list", "std"), "go list std")), "\n"), "C")
	}
	stdMutex.Unlock()

	if tagsCustom != "" {
		tags = append(tags, strings.Split(tagsCustom, " ")...)
	}

	importMap := make(map[string]struct{})

	imp := "Deps"
	if onlyDirect {
		imp = "Imports"
	}

	//TODO: cache
	cmd := utils.GoList("'{{ join .TestImports \"|\" }}':'{{ join .XTestImports \"|\" }}':'{{ join ."+imp+" \"|\" }}'", fmt.Sprintf("-tags=\"%v\"", strings.Join(tags, "\" \"")))
	cmd.Dir = path
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", k, v))
	}

	out := strings.TrimSpace(utils.RunCmd(cmd, "go list deps"))
	out = strings.Replace(out, "'", "", -1)
	out = strings.Replace(out, ":", "|", -1)
	libs := strings.Split(out, "|")

	for i := len(libs) - 1; i >= 0; i-- {
		if strings.TrimSpace(libs[i]) == "" {
			libs = append(libs[:i], libs[i+1:]...)
			continue
		}

		importedMutex.Lock()
		dep, ok := imported[libs[i]]
		importedMutex.Unlock()
		if ok {
			importMap[dep] = struct{}{}
			libs = append(libs[:i], libs[i+1:]...)
			continue
		}

		for _, k := range std {
			if libs[i] == k {
				libs = append(libs[:i], libs[i+1:]...)
				break
			}
		}
	}

	wg := new(sync.WaitGroup)
	wc := make(chan bool, 50)
	wg.Add(len(libs))
	for _, l := range libs {
		wc <- true
		go func(l string) {
			defer func() {
				<-wc
				wg.Done()
			}()

			if strings.Contains(l, "github.com/peterq/pan-light/qt") && !strings.Contains(l, "qt/tool-chain") {
				return
			}

			cmd := utils.GoList("{{.Dir}}", fmt.Sprintf("-tags=\"%v\"", strings.Join(tags, "\" \"")), l)
			for k, v := range env {
				cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", k, v))
			}

			dep := strings.TrimSpace(utils.RunCmd(cmd, "go list dir"))
			if dep == "" {
				return
			}

			importedMutex.Lock()
			importMap[dep] = struct{}{}
			imported[l] = dep
			importedMutex.Unlock()
		}(l)
	}
	wg.Wait()

	var imports []string
	for k := range importMap {
		imports = append(imports, k)
	}
	return imports
}

func GetGoFiles(path, target, tagsCustom string) []string {
	utils.Log.WithField("path", path).WithField("target", target).WithField("tagsCustom", tagsCustom).Debug("get go files")

	env, tags, _, _ := BuildEnv(target, "", "")
	if tagsCustom != "" {
		tags = append(tags, strings.Split(tagsCustom, " ")...)
	}

	//TODO: cache
	cmd := utils.GoList("'{{ join .GoFiles \"|\" }}':'{{ join .CgoFiles \"|\" }}':'{{ join .TestGoFiles \"|\" }}':'{{ join .XTestGoFiles \"|\" }}'", fmt.Sprintf("-tags=\"%v\"", strings.Join(tags, "\" \"")))
	cmd.Dir = path
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", k, v))
	}

	out := strings.TrimSpace(utils.RunCmd(cmd, "go list gofiles"))
	out = strings.Replace(out, "'", "", -1)
	out = strings.Replace(out, ":", "|", -1)

	importMap := make(map[string]struct{})
	for _, v := range strings.Split(out, "|") {
		if strings.TrimSpace(v) != "" {
			importMap[v] = struct{}{}
		}
	}

	olibs := make([]string, 0)
	for k := range importMap {
		olibs = append(olibs, filepath.Join(path, k))
	}
	return olibs
}
