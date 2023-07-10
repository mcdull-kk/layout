package env

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var envCache struct {
	once sync.Once
	m    map[string]string
}

func Getenv(key string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	envCache.once.Do(initEnvCache)
	return envCache.m[key]
}

func initEnvCache() {
	envCache.m = make(map[string]string)
	file, _ := envFile()
	if file == "" {
		return
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	for len(data) > 0 {
		// Get next line.
		line := data
		i := bytes.IndexByte(data, '\n')
		if i >= 0 {
			line, data = line[:i], data[i+1:]
		} else {
			data = nil
		}

		i = bytes.IndexByte(line, '=')
		if i < 0 || line[0] < 'A' || 'Z' < line[0] {
			// Line is missing = (or empty) or a comment or not a valid env name. Ignore.
			// (This should not happen, since the file should be maintained almost
			// exclusively by "go env -w", but better to silently ignore than to make
			// the go command unusable just because somehow the env file has
			// gotten corrupted.)
			continue
		}
		key, val := line[:i], line[i+1:]
		envCache.m[string(key)] = string(val)
	}
}

// EnvFile returns the name of the Go environment configuration file.
func envFile() (string, error) {
	if file := os.Getenv("GOENV"); file != "" {
		if file == "off" {
			return "", fmt.Errorf("GOENV=off")
		}
		return file, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("missing user-config dir")
	}
	return filepath.Join(dir, "go/env"), nil
}

func Gopath() (gp string) {
	gopaths := strings.Split(Getenv("GOPATH"), string(filepath.ListSeparator))

	if len(gopaths) == 1 && gopaths[0] != "" {
		return gopaths[0]
	}
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	abspwd, err := filepath.Abs(pwd)
	if err != nil {
		return
	}
	for _, gopath := range gopaths {
		if gopath == "" {
			continue
		}
		absgp, err := filepath.Abs(gopath)
		if err != nil {
			return
		}
		if strings.HasPrefix(abspwd, absgp) {
			return absgp
		}
	}
	return build.Default.GOPATH
}
