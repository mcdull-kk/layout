package main

import (
	"fmt"
	"os"

	"github.com/mcdull-kk/layout/tool"
	"github.com/mcdull-kk/layout/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "mcdull"
	app.Usage = "mcdull toolkit"
	app.Version = version.Version
	app.Commands = []*cli.Command{
		{
			Name:            "new",
			Aliases:         []string{"n", "c"},
			Usage:           "create a new project",
			Action:          tool.GenNewProjectAction,
			SkipFlagParsing: true,
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "build",
			Action:  tool.BuildAction,
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run",
			Action:  tool.RunAction,
		},
		{
			Name:            "tool",
			Aliases:         []string{"t"},
			Usage:           "tool",
			Action:          tool.ToolKitAction,
			SkipFlagParsing: true,
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version",
			Action: func(c *cli.Context) error {
				fmt.Println(app.Version)
				return nil
			},
		},
		{
			Name:    "upgrade",
			Aliases: []string{"u"},
			Usage:   "upgrade",
			Action:  tool.UpgradeAction,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
