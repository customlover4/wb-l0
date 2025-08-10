package templates

import (
	"path/filepath"
	"runtime"
	"text/template"
)

func GetTemplatesPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "templates")
}

func LoadTemplates() (*template.Template, error) {
	tplPath := GetTemplatesPath()
	return template.ParseGlob(filepath.Join(tplPath, "**/*.html"))
}
