package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/mcdull-kk/layout/env"
)

var (
	provideSupportProtocs = []*Proctoc{
		{
			_name:   "bm",
			_gen:    "protoc-gen-bm",
			_getgen: "go get -u github.com/mcdull-kk/layout/toolkit/protobuf/protoc-gen-bm",
			_protoc: "protoc --proto_path=%s --proto_path=%s --proto_path=%s --bm_out=:.",
		},
		{
			_name:   "ecode",
			_gen:    "protoc-gen-ecode",
			_getgen: "go get -u github.com/mcdull-kk/layout/toolkit/protobuf/protoc-gen-ecode",
			_protoc: "protoc --proto_path=%s --proto_path=%s --proto_path=%s --ecode_out=" +
				"Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:.",
		},
		{
			_name:   "grpc",
			_gen:    "protoc-gen-gofast", // protoc-gen-gofast 性能高于 protoc-gen-go
			_getgen: "go get -u github.com/gogo/protobuf/protoc-gen-gofast",
			_protoc: "protoc --proto_path=%s --proto_path=%s --proto_path=%s --gofast_out=plugins=grpc," +
				"Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:.",
		},
		{
			_name:   "grpc-gateway",
			_gen:    "protoc-gen-grpc-gateway",
			_getgen: "go get github.com/grpc-ecosystem/grpc-gateway",
			_protoc: "protoc --proto_path=%s --proto_path=%s --proto_path=%s --grpc-gateway_out=" +
				"Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types," +
				"Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:.",
		},
	}
)

type Proctoc struct {
	_name   string
	_gen    string
	_getgen string
	_protoc string
}

func (p *Proctoc) installGen() error {
	if _, err := exec.LookPath(p._gen); err != nil {
		if err := p.goget(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Proctoc) gen(files []string) error {
	return p.generate(files)
}

func (p *Proctoc) goget() error {
	args := strings.Split(p._getgen, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println(p._gen)
	return cmd.Run()
}

func (p *Proctoc) generate(files []string) error {
	pwd, _ := os.Getwd()
	gosrc := path.Join(env.Gopath(), "src")
	ext, err := p.latest()
	if err != nil {
		return err
	}
	line := fmt.Sprintf(p._protoc, gosrc, ext, pwd)
	log.Println(line, strings.Join(files, " "))
	args := strings.Split(line, " ")
	args = append(args, files...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = pwd
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *Proctoc) latest() (string, error) {
	gopath := env.Gopath()
	ext := path.Join(gopath, "src/github.com/mcdull-kk/layout/third_party")
	if _, err := os.Stat(ext); !os.IsNotExist(err) {
		return ext, nil
	}
	ext = path.Join(gopath, "src/mcdull/third_party")
	if _, err := os.Stat(ext); !os.IsNotExist(err) {
		return ext, nil
	}
	baseMod := path.Join(gopath, "pkg/mod/github.com/mcdull-kk")
	files, err := ioutil.ReadDir(baseMod)
	if err != nil {
		return "", err
	}
	for i := len(files) - 1; i >= 0; i-- {
		if strings.HasPrefix(files[i].Name(), "layout@") {
			return path.Join(baseMod, files[i].Name(), "third_party"), nil
		}
	}
	return "", errors.New("not found layout package")
}
