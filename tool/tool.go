package tool

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mcdull-kk/layout/env"
	"github.com/mcdull-kk/layout/version"
)

var (
	_gname     string
	GenProject string = "genproject"
)

func init() {
	flag.StringVar(&_gname, "project", "mcdull", "prject name")
}

var (
	provideSupportTools = []*ToolKit{
		{
			Name:      "mcdull",
			Alias:     "mcdull",
			BuildTime: time.Date(2023, 7, 8, 0, 0, 0, 0, time.Local),
			Install:   "go get -u github.com/mcdull-kk/layout@" + version.Version,
			Summary:   "mcdull toolkit",
			Platform:  []string{"darwin", "linux", "windows"},
			Author:    "mcdull",
			Hidden:    true,
		},
		{
			Name:      GenProject,
			Alias:     fmt.Sprintf("%s-gen-project", _gname),
			BuildTime: time.Date(2023, 7, 8, 0, 0, 0, 0, time.Local),
			Install:   "go get -u github.com/mcdull-kk/layout/toolkit/gen-project@" + version.Version,
			Platform:  []string{"darwin", "linux", "windows"},
			Author:    "mcdull",
			Hidden:    true,
		},
	}
)

func getToolKitList() []*ToolKit {
	return provideSupportTools
}

type ToolKit struct {
	Name         string    `json:"name"`
	Alias        string    `json:"alias"`
	BuildTime    time.Time `json:"build_time"`
	Install      string    `json:"install"`
	Requirements []string  `json:"requirements"`
	Dir          string    `json:"dir"`
	Summary      string    `json:"summary"`
	Platform     []string  `json:"platform"`
	Author       string    `json:"author"`
	URL          string    `json:"url"`
	Hidden       bool      `json:"hidden"`
	requires     []*ToolKit
}

func (t ToolKit) installed() bool {
	_, err := os.Stat(t.path())
	return err == nil
}

func (t ToolKit) path() string {
	name := t.Alias
	if name == "" {
		name = t.Name
	}
	gobin := env.Getenv("GOBIN")
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if gobin != "" {
		return filepath.Join(gobin, name)
	}
	return filepath.Join(env.Gopath(), "bin", name)
}

func (t ToolKit) needUpdated() bool {
	for _, r := range t.requires {
		if r.needUpdated() {
			return true
		}
	}
	if !t.supportOS() || t.Install == "" {
		return false
	}
	if f, err := os.Stat(t.path()); err == nil {
		if t.BuildTime.After(f.ModTime()) {
			return true
		}
	}
	return false
}

func (t ToolKit) supportOS() bool {
	for _, p := range t.Platform {
		if strings.ToLower(p) == runtime.GOOS {
			return true
		}
	}
	return false
}

func (t ToolKit) install() {
	if t.Install == "" {
		fmt.Fprint(os.Stderr, color.RedString("%s: install is empty\n", t.Name))
		return
	}
	fmt.Println(t.Install)
	cmds := strings.Split(t.Install, " ")
	if len(cmds) > 0 {
		if err := runTool(t.Name, path.Dir(t.path()), cmds[0], cmds[1:]); err == nil {
			color.Green("%s: install success!", t.Name)
		}
	}
}
