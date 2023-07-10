package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	withBM         bool
	withGRPC       bool
	withGrpcGatewy bool
	withSwagger    bool
	withEcode      bool
)

func protocAction(ctx *cli.Context) (err error) {
	if err = checkProtoc(); err != nil {
		return err
	}
	files := ctx.Args().Slice()
	if len(files) == 0 {
		files, _ = filepath.Glob("*.proto")
	}
	if !withGRPC && !withBM && !withSwagger && !withEcode {
		withBM = true
		withGRPC = true
		withSwagger = true
		withEcode = true
	}
	if withBM {
		if err = installBMGen(); err != nil {
			return
		}
		if err = genBM(files); err != nil {
			return
		}
	}
	if withGRPC {
		if err = installGRPCGen(); err != nil {
			return err
		}
		if err = genGRPC(files); err != nil {
			return
		}
	}
	if withSwagger {
		if err = installSwaggerGen(); err != nil {
			return
		}
		if err = genSwagger(files); err != nil {
			return
		}
	}
	if withEcode {
		if err = installEcodeGen(); err != nil {
			return
		}
		if err = genEcode(files); err != nil {
			return
		}
	}
	if withGrpcGatewy {

	}
	log.Printf("generate %s success.\n", strings.Join(files, " "))
	return nil
}

func checkProtoc() error {
	if _, err := exec.LookPath("protoc"); err != nil {
		switch runtime.GOOS {
		case "darwin":
			fmt.Println("brew install protobuf")
			cmd := exec.Command("brew", "install", "protobuf")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				return err
			}
		case "linux":
			fmt.Println("snap install --classic protobuf")
			cmd := exec.Command("snap", "install", "--classic", "protobuf")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				return err
			}
		default:
			return errors.New("you haven't installed protobuf yet,please install it manually：https://github.com/protocolbuffers/protobuf/releases")
		}
	}
	return nil
}
