package version

import (
	"bytes"
	"html/template"
	"runtime"
)

var (
	Version = "v1.0.0"
)

type VersionOptions struct {
	GitCommit string
	Version   string
	GoVersion string
	Os        string
	Arch      string
}

var _versionTemplate = ` Version:      {{.Version}}
 Go version:   {{.GoVersion}}
 OS/Arch:      {{.Os}}/{{.Arch}}
 `

func GetVersion() string {
	var doc bytes.Buffer
	vo := VersionOptions{
		Version:   Version,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
	tmpl, _ := template.New("version").Parse(_versionTemplate)
	tmpl.Execute(&doc, vo)
	return doc.String()
}
