package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/peterq/pan-light/qt/tool-chain/binding/templater"

	"github.com/peterq/pan-light/qt/tool-chain/cmd/deploy"
	"github.com/peterq/pan-light/qt/tool-chain/cmd/minimal"
	"github.com/peterq/pan-light/qt/tool-chain/cmd/moc"

	"github.com/peterq/pan-light/qt/tool-chain/utils"
)

func Test(target string, docker, vagrant bool, vagrantsystem string) {
	if docker && target == "darwin" {
		utils.Log.Warn("darwin is currently not supported as a deploy target with docker; testing the linux deployment instead")
		target = "linux"
	}

	utils.Log.Infof("running: 'qtsetup test %v' [docker=%v] [vagrant=%v]", target, docker, vagrant)

	if utils.CI() && target == runtime.GOOS && runtime.GOOS != "windows" { //TODO: test on windows
		utils.Log.Infof("running setup/test %v CI", target)

		path := utils.GoQtPkgPath("tool-chain", "cmd", "moc", "test")

		moc.Moc(path, target, "", false, false)
		minimal.Minimal(path, target, "")

		var pattern string
		if strings.Contains(runtime.Version(), "1.1") || strings.Contains(runtime.Version(), "devel") {
			pattern = "all="
		}

		cmd := exec.Command("go", "test", "-v", "-tags=minimal", fmt.Sprintf("-ldflags=%v\"-s\"", pattern))
		cmd.Env = append(os.Environ(), "GODEBUG=cgocheck=2")
		cmd.Dir = path
		utils.RunCmd(cmd, "run \"go test\"")
	}

	mode := "test"
	var examples map[string][]string
	if utils.CI() {
		mode = "build"
		examples = map[string][]string{
			"androidextras": {"jni", "notification"},

			"canvas3d": {"framebuffer", "interaction", "jsonmodels",
				"quickitemtexture", "textureandlight",
				filepath.Join("threejs", "cellphone"),
				filepath.Join("threejs", "oneqt"),
				filepath.Join("threejs", "planets"),
			},

			"charts": {"audio"},

			"common": {"qml_demo", "widgets_demo"},

			//"grpc": []string{"hello_world","hello_world2"},

			//"gui": []string{"analogclock", "openglwindow", "rasterwindow"},

			//opengl: []string{"2dpainting"},

			"qml": {"adding", "application", "drawer_nav_x",
				filepath.Join("extending", "chapter1-basics"),
				filepath.Join("extending", "chapter2-methods"),
				filepath.Join("extending", "chapter3-bindings"),
				filepath.Join("extending", "chapter4-customPropertyTypes"),
				filepath.Join("extending", "components", "test_dir"),
				filepath.Join("extending", "components", "test_dir_2"),
				filepath.Join("extending", "components", "test_go"),
				filepath.Join("extending", "components", "test_module"),
				filepath.Join("extending", "components", "test_module_2"),
				filepath.Join("extending", "components", "test_qml"),
				filepath.Join("extending", "components", "test_qml_go"),
				"gallery", "material",
				//filepath.Join("printslides", "cmd", "printslides"),
				"prop", "prop2" /*"quickflux", "webview"*/},

			"qt3d": {"audio-visualizer-qml"},

			"quick": {"bridge", "bridge2", "calc", "dialog", "dynamic",
				"hotreload", "listview", "sailfish", "tableview", "translate", "view"},

			"sailfish": {"listview", "listview_variant"},

			"showcases": {"sia"},

			"sql": {"masterdetail", "masterdetail_qml", "querymodel"},

			"uitools": {"calculator"},

			"webchannel": {"chatserver-go" /*"standalone" "webview"*/},

			"widgets": {"bridge2" /*"dropsite"*/, "graphicsscene", "line_edits", "pixel_editor",
				/*"renderer"*/ "share", "systray" /*"table"*/, "textedit", filepath.Join("treeview", "treeview_dual"),
				filepath.Join("treeview", "treeview_filelist"), "video_player" /*"webengine"*/, "xkcd"},
		}
	} else {
		if strings.HasPrefix(target, "sailfish") {
			examples = map[string][]string{
				"quick": {"sailfish"},

				"sailfish": {"listview", "listview_variant"},
			}
		} else {
			examples = map[string][]string{
				"qml": {"application", "drawer_nav_x", "gallery"},

				"quick": {"calc"},

				"widgets": {"line_edits", "pixel_editor", "textedit"},
			}
		}
	}

	if utils.QT_VAGRANT() && target == "ios-simulator" {
		mode = "build"
	}

	for cat, list := range examples {
		for _, example := range list {
			if target != runtime.GOOS && example == "textedit" {
				continue
			}

			if (target == "js" || target == "wasm") &&
				cat == "charts" || cat == "uitools" || cat == "sql" ||
				cat == "androidextras" || cat == "qt3d" || cat == "webchannel" ||
				(cat == "widgets" && strings.HasPrefix(example, "treeview")) ||
				example == "video_player" {
				continue
			}

			example := filepath.Join(cat, example)

			path := filepath.Join(strings.TrimSpace(utils.RunCmdOptional(utils.GoList("{{.Dir}}", "github.com/peterq/pan-light/qt/tool-chain/examples"), "get doc dir")), example)
			utils.Log.Infof("testing %v", example)
			deploy.Deploy(
				mode,
				target,
				path,
				docker,
				"",
				"",
				false,
				"",
				vagrant,
				vagrantsystem,
				false,
			)
			templater.CleanupDepsForCI()
			templater.CleanupDepsForCI = func() {}
		}
	}
}
