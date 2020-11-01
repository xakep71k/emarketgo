package magazine

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"text/template"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

func loadDataFromJsonFile(dataPath string, v interface{}) error {
	data, err := ioutil.ReadFile(dataPath)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}

func doMinify(body []byte, mtype string) ([]byte, error) {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	b, err := m.Bytes(mtype, body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func execute(name, templateStr string, input interface{}) (*bytes.Buffer, error) {
	templ, err := template.New(name).Funcs(defaultTemplateFuncs()).Parse(templateStr)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := templ.Execute(buf, input); err != nil {
		return nil, err
	}

	return buf, nil
}
