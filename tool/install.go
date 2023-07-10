package tool

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
)

func runToolKit(name string, args []string) (err error) {
	for _, t := range getToolKitList() {
		if name == t.Name {
			if !t.installed() || t.needUpdated() {
				t.install()
			}
			pwd, _ := os.Getwd()
			err = runTool(t.Name, pwd, t.path(), args)
			return
		}
	}
	return fmt.Errorf("not found %s", name)
}

func runTool(name, dir, cmd string, args []string) (err error) {
	toolCmd := &exec.Cmd{
		Path:   cmd,
		Args:   append([]string{cmd}, args...),
		Dir:    dir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}
	if filepath.Base(cmd) == cmd {
		var lp string
		if lp, err = exec.LookPath(cmd); err == nil {
			toolCmd.Path = lp
		}
	}
	if err = toolCmd.Run(); err != nil {
		if e, ok := err.(*exec.ExitError); !ok || !e.Exited() {
			fmt.Fprintf(os.Stderr, "run %s error: %v\n", name, err)
		}
	}
	return
}

func checkInstall(name string) (err error) {
	for _, t := range getToolKitList() {
		if name == t.Name {
			if !t.installed() || t.needUpdated() {
				t.install()
			}
			return
		}
	}
	return fmt.Errorf("not found %s", name)
}

func install(name string) {
	if name == "" {
		fmt.Fprint(os.Stderr, color.HiRedString("please write the name of the tool you want to install\n"))
		return
	}
	for _, t := range getToolKitList() {
		if name == t.Name {
			t.install()
			return
		}
	}
	fmt.Fprint(os.Stderr, color.HiRedString("install fail not found %s\n", name))
}

func installAll() {
	for _, t := range getToolKitList() {
		if t.Install != "" {
			t.install()
		}
	}
}

func upgradeAll() {
	for _, t := range getToolKitList() {
		if t.needUpdated() {
			t.install()
		}
	}
}

func getNotice(t *ToolKit) (notice string) {
	if !t.supportOS() || t.Install == "" {
		return
	}
	notice = color.HiGreenString("(uninstalled)")
	if t.installed() {
		notice = color.HiBlueString("(installed)")
		if t.needUpdated() {
			notice = color.RedString("(needUpdate)")
		}
	}
	return
}
