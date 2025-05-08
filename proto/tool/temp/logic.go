package temp

import (
	"bytes"
	"embed"
	"gmicro/pkg/gerr"
	"gmicro/pkg/log"
	"text/template"
)

//go:embed gtpl
var tplFs embed.FS

func GenCodeByTemplate(name string, data interface{}) (string, error) {
	t, err := LoadTemplate(name, nil)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	var b bytes.Buffer
	err = t.Execute(&b, data)
	if err != nil {
		log.Errorf("execute temp err:%v", err)
		return "", err
	}

	return b.String(), nil
}

func LoadTemplate(name string, funcMap template.FuncMap) (*template.Template, error) {
	buf, err := tplFs.ReadFile("gtpl/" + name)
	if err != nil {
		return nil, gerr.Wrap(err)
	}

	var t *template.Template
	if funcMap != nil && len(funcMap) > 0 {
		t, err = template.New(name).Funcs(funcMap).Parse(string(buf))
	} else {
		t, err = template.New(name).Parse(string(buf))
	}
	return t, err
}
