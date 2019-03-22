package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
		} {
			fmt.Printf("  %v%v%v\n", m.name, strings.Repeat(" ", 12-len(m.name)), m.desc)
		}
		println()

		os.Exit(0)
	}
	flag.Parse()

	cmd := "pc"

	cmd = flag.Arg(0)

	switch cmd {
	case "pc":
		pcCmd()
	default:
		flag.Usage()
	}
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
		runCmd("./pc/gui", "go", "build", "-tags=plugin", "--buildmode=plugin", "-o", "gui-plugin.so", "gui-plugin.go")
	}
	log.Println("打包qml...")
	cmd(qtBin("rcc"), "-binary", "pc/gui/qml/qml.qrc", "-o", "pc/gui/qml/qml.rcc").Run()
	log.Println("启动客户端...")
	runCmd("./pc", "go", "run", "-tags=plugin", "pan-light-pc-dev.go")
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
	c.Run()
}
