package tool

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/mcdull-kk/layout/version"
	"github.com/urfave/cli/v2"
)

func GenNewProjectAction(c *cli.Context) (err error) {
	return runToolKit(GenProject, c.Args().Slice())
}

func ToolKitAction(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		list := getToolKitList()
		sort.Slice(list, func(i, j int) bool { return list[i].BuildTime.After(list[j].BuildTime) })
		for _, t := range list {
			if t.Hidden {
				continue
			}
			updateTime := t.BuildTime.Format("2006/01/02")
			fmt.Printf("%s%s: %s Author(%s) [%s]\n", color.HiMagentaString(t.Name), getNotice(t), color.HiCyanString(t.Summary), t.Author, updateTime)
		}
		fmt.Println("\ninistall tool : mcdull tool install demo")
		fmt.Println("exec tool: mcdull tool demo")
		fmt.Println("install all toolkit: mcdull tool install all")
		fmt.Println("upgrade all: mcdull tool upgrade all")
		return
	}
	commond := c.Args().First()
	switch commond {
	case "upgrade":
		upgradeAll()
		return
	case "install":
		name := c.Args().Get(1)
		if name == "all" {
			installAll()
		} else {
			install(name)
		}
		return
	case "check_install":
		if e := checkInstall(c.Args().Get(1)); e != nil {
			fmt.Fprintf(os.Stderr, "%v\n", e)
		}
		return
	}
	if e := runToolKit(commond, c.Args().Slice()[1:]); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
	}
	return
}

func UpgradeAction(c *cli.Context) (err error) {
	install("mcdull")
	return nil
}

func BuildAction(c *cli.Context) error {
	base, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	args := append([]string{"build"}, c.Args().Slice()...)
	cmd := exec.Command("go", args...)
	cmd.Dir = buildDir(base, "cmd", 5)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("directory: %s\n", cmd.Dir)
	fmt.Printf("mcdull: %s\n", version.Version)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Println("build success.")
	return nil
}

func RunAction(c *cli.Context) error {
	base, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir := buildDir(base, "cmd", 5)
	conf := path.Join(filepath.Dir(dir), "configs")
	args := append([]string{"run", "main.go", "-conf", conf}, c.Args().Slice()...)
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return nil
}

func buildDir(base string, cmd string, n int) string {
	dirs, err := ioutil.ReadDir(base)
	if err != nil {
		panic(err)
	}
	for _, d := range dirs {
		if d.IsDir() && d.Name() == cmd {
			return path.Join(base, cmd)
		}
	}
	if n <= 1 {
		return base
	}
	return buildDir(filepath.Dir(base), cmd, n-1)
}
