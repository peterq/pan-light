package main

// 这个是为了方便开发写的脚本, 代码写的很糙, 见笑

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	flag.Usage = func() {
		println("用法: pan-light [cmd] [sub cmd]\n")

		println("参数:\n")
		flag.PrintDefaults()
		println()

		println("命令:\n")
		for _, m := range []struct{ name, desc string }{
			{"pc dev", "启动pc客户端, 开发模式, 用plugin加速编译, 不支持windows"},
			{"pc start", "启动pc客户端"},
			{"pc moc", "生成moc"},
			{"server start", "启动server"},
		} {
			fmt.Printf("  %v%v%v\n", m.name, strings.Repeat(" ", 12-len(m.name)), m.desc)
		}
		println()

		os.Exit(0)
	}
	flag.Parse()

	os.Setenv("GOPROXY", "https://goproxy.io")

	cmd := "pc"

	cmd = flag.Arg(0)

	switch cmd {
	case "pc":
		pcCmd()
	case "demo":
		demoCmd()
	case "server":
		serverCmd()
	default:
		flag.Usage()
	}
}

func serverCmd() {
	cmd := flag.Arg(1)
	switch cmd {
	case "start":
		serverStart()
	default:
		flag.Usage()
	}
}

func serverStart() {
	os.Setenv("pan_light_server_conf", "pan-light-server.yaml")
	runCmd("./server", "go", "run", "pan-light-server.go")
}

func demoCmd() {
	cmd := flag.Arg(1)
	switch cmd {
	case "host":
		demoHost()
	case "build.docker.slave":
		buildDockerSlave()
	case "avatar":
		demoAvatar()
	default:
		flag.Usage()
	}
}

func demoAvatar() {
	runCmd("./server/cmd", "go", "run", "main.go", "avatar")
}

func buildDockerSlave() {
	runCmd("demo", "go", "build", "-o", "./slave/ubuntu16.04/root.pan-light/demo_instance_manager", "slave.go")
	//runCmd("demo", "docker", "build", "-t", "pan-light-slave", "./slave/ubuntu16.04")
}

func demoHost() {
	cmd(qtBin("rcc"), "-binary", "pc/gui/qml/qml.qrc", "-o", "demo/slave/ubuntu16.04/root.pan-light/lib/qml.rcc").Run()
	runCmd("./demo", "go", "build", "-o", "./slave/ubuntu16.04/root.pan-light/demo_instance_manager", "slave.go")
	runCmd("./demo", "go", "run", "host.go")
}

func pcCmd() {
	cmd := flag.Arg(1)
	switch cmd {
	case "start":
		pcStart()
	case "moc":
		pcMoc()
	case "dev":
		pcDev()
	case "build":
		pcBuild()
	case "download-icon":
		downloadIcon()
	default:
		flag.Usage()
	}
}

func pcDev() {
	var rebuildPlugin bool
	flag.BoolVar(&rebuildPlugin, "rebuild", false, "重新编译plugin插件")
	flag.Parse()
	pluginPath := "./pc/gui/gui-plugin.so"
	_, err := os.Stat(pluginPath)
	if os.IsNotExist(err) || rebuildPlugin {
		log.Println("编译gui插件...")
		runCmd("./pc", "go", "build", "-p", strconv.Itoa(runtime.NumCPU()-1), "-v", "-tags=plugin", "--buildmode=plugin", "-o", "gui/gui-plugin.so", "gui/gui-plugin.go")
	}
	log.Println("打包qml...")
	cmd(qtBin("rcc"), "-binary", "pc/gui/qml/qml.qrc", "-o", "pc/gui/qml/qml.rcc").Run()
	log.Println("启动客户端...")
	runCmd("./pc", "go", "run", "-tags=plugin", "pan-light-pc-dev.go")
}

func pcBuild() {
	runCmd("./pc", "go", "run", "../qt/cmd/qtdeploy/main.go", "-fast", "build", "desktop")
	// 正式包编译步骤
	// cd pc
	// go run ../qt/cmd/qtrcc/main.go desktop gui/qml
	// go build -p 8 -v -ldflags=all="-s" -o deploy/linux/pc -tags="production"

}

func qtBin(name string) string {
	v, ok1 := os.LookupEnv("QT_VERSION")
	d, ok2 := os.LookupEnv("QT_DIR")
	if !ok1 || !ok2 {
		panic("请先配置qt环境变量")
	}
	return d + "/" + v + "/gcc_64/bin/" + name
}

func pcStart() {
	// 打包qml
	log.Println("打包qml...")
	c := cmd("go", "run", "../qt/cmd/qtrcc/main.go", "desktop", "gui/qml")
	c.Dir, _ = filepath.Abs("./pc")
	c.Run()

	// 启动
	log.Println("启动pc客户端")
	c = cmd("go", "run", "pan-light-pc.go")
	c.Dir = "./pc"
	c.Run()
}

// 清除svg的style font节点, 减小尺寸
func iconSimplify() {
	dir := "./pc/gui/qml/assets/images/icons"
	ls, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(`\<defs\>(.|\n)*\</defs\>`)
	for _, f := range ls {
		if strings.Index(f.Name(), ".svg") == len(f.Name())-4 {
			filename := dir + "/" + f.Name()
			bin, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			content := string(bin)
			target := reg.ReplaceAllString(content, "")
			if target != content {
				err = ioutil.WriteFile(filename, []byte(target), os.ModePerm)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func pcMoc() {
	log.Println("moc...")
	c := cmd("go", "run", "../qt/cmd/qtmoc/main.go", "desktop", "gui/comp")
	c.Dir, _ = filepath.Abs("./pc")
	e := c.Run()
	if e != nil {
		log.Println(e)
	}
}

func cmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd
}

func runCmd(path, name string, arg ...string) {
	c := cmd(name, arg...)
	c.Dir, _ = filepath.Abs(path)
	c.Env = append(os.Environ(), envFile(c.Dir)...)
	c.Run()
}

func envFile(path string) (env []string) {
	bin, err := ioutil.ReadFile(path + "/.env")
	if err != nil {
		return
	}
	str := string(bin)
	for _, ln := range strings.Split(str, "\n") {
		ln = strings.Trim(ln, "\n")
		if ln == "" {
			continue
		}
		if ln[0] == '#' {
			continue
		}
		env = append(env, ln)
	}
	return
}

type gson = map[string]interface{}
type binary = []byte

func downloadIcon() {
	u := "https://www.iconfont.cn/api/collection/detail.json?id=2271"
	r, e := http.Get(u)
	panicIf(e)
	bin, e := ioutil.ReadAll(r.Body)
	panicIf(e)
	var j gson
	panicIf(json.Unmarshal(bin, &j))
	reg := regexp.MustCompile("[a-z0-9]+$")
	var imgs []string
	for _, icon := range j["data"].(gson)["icons"].([]interface{}) {
		s := icon.(gson)["show_svg"].(string)
		name := icon.(gson)["name"].(string)
		name = string(reg.Find(binary(name)))
		if name == "" {
			name = icon.(gson)["name"].(string)
			if strings.Contains(name, "未知") {
				name = "unknown"
			} else {
				continue
			}
		}
		imgs = append(imgs, name)
		e = ioutil.WriteFile("pc/gui/qml/assets/images/icons/file/"+name+".svg", binary(s), os.ModePerm)
		panicIf(e)
	}
	bin, _ = json.Marshal(imgs)
	log.Println(string(bin))
}

func panicIf(e error) {
	if e != nil {
		panic(e)
	}

}
