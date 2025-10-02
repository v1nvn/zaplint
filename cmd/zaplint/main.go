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
	v, vcsInfo := getVersion()
	fmt.Printf("zaplint version %s %s/%s (built with %s)\n", v, runtime.GOOS, runtime.GOARCH, runtime.Version())
	if vcsInfo != "" {
		fmt.Println(vcsInfo)
	}
	os.Exit(0)
	return nil
}

func getVersion() (string, string) {
	var vcsInfo string

	if version != "dev" {
		return version, ""
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		} else {
			version = "dev"
		}

		// Extract VCS information
		var revision, time, modified string
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.time":
				time = setting.Value
			case "vcs.modified":
				modified = setting.Value
			}
		}

		if revision != "" {
			vcsInfo = fmt.Sprintf("commit: %s", revision)
			if time != "" {
				vcsInfo += fmt.Sprintf(", built at: %s", time)
			}
			if modified == "true" {
				vcsInfo += " (dirty)"
			}
		}

		return version, vcsInfo
	}

	return "dev", ""
}
