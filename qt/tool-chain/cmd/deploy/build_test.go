package deploy

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_escapeFlags(t *testing.T) {

	tmpfile := filepath.Join(os.TempDir(), "escapeFlags.go")
	defer os.Remove(tmpfile)
	if err := ioutil.WriteFile(tmpfile, []byte("package main;var foo, abc string;func main() { println(foo, abc) }"), 0644); err != nil {
		t.Fatal(err)
	}

	var pattern string
	if strings.Contains(runtime.Version(), "1.1") || strings.Contains(runtime.Version(), "devel") {
		pattern = "all="
	}

	for _, flags := range [][]string{
		{},
		{"-w"},
		{"-w", "-s"},
		{"-w", "-s", "-extldflags=-v"},
		{"-w", "-s", "-extldflags=-v"},
	} {
		for i, tc := range []string{
			"",
			"-X main.foo=bar",
			"-X \"main.foo=bar\"",
			"-X 'main.foo=bar'",
		} {
			cmd := exec.Command("go", "run", "-v", fmt.Sprintf("-ldflags=%v%v", pattern, escapeFlags(flags, tc)), tmpfile)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err, string(out), cmd.Args)
			}
			outC := "bar"
			if tc == "" {
				outC = ""
			}
			if outT := strings.TrimSpace(string(out)); outT != outC {
				t.Fatal(i, outT, "!=", outC)
			}
		}

		for i, tc := range []string{
			"",
			"-X \"main.foo=bar baz\" -X \"main.abc=bbb ddd\"",
			"-X \"main.foo=bar baz\" -X 'main.abc=bbb ddd'",
			"-X 'main.foo=bar baz' -X 'main.abc=bbb ddd'",
		} {
			cmd := exec.Command("go", "run", "-v", fmt.Sprintf("-ldflags=%v%v", pattern, escapeFlags(flags, tc)), tmpfile)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err, string(out), cmd.Args)
			}
			outC := "bar baz bbb ddd"
			if tc == "" {
				outC = ""
			}
			if outT := strings.TrimSpace(string(out)); outT != outC {
				t.Fatal(i, outT, "!=", outC)
			}
		}
	}
}
