package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/v1nvn/zaplint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var version = "dev"

func main() {
	// override the builtin -V flag.
	flag.Var(versionFlag{}, "V", "print version and exit")
	singlechecker.Main(zaplint.New(nil))
}

type versionFlag struct{}

func (versionFlag) String() string   { return "" }
func (versionFlag) IsBoolFlag() bool { return true }
func (versionFlag) Set(string) error {
	v := getVersion()
	fmt.Printf("zaplint version %s %s/%s (built with %s)\n", v, runtime.GOOS, runtime.GOARCH, runtime.Version())
	os.Exit(0)
	return nil
}

func getVersion() string {
	if version != "dev" {
		return version
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	return "dev"
}
