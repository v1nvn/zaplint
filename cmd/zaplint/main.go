package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/v1nvn/zaplint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var version = "dev" // injected at build time.

func main() {
	// override the builtin -V flag.
	flag.Var(versionFlag{}, "V", "print version and exit")
	singlechecker.Main(zaplint.New(nil))
}

type versionFlag struct{}

func (versionFlag) String() string   { return "" }
func (versionFlag) IsBoolFlag() bool { return true }
func (versionFlag) Set(string) error {
	fmt.Printf("zaplint version %s %s/%s (built with %s)\n", version, runtime.GOOS, runtime.GOARCH, runtime.Version())
	os.Exit(0)
	return nil
}
