package template

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"
)

type Template struct {
	template *template.Template

	// Global template directory
	templatePath string

	// Global template functions
	templateFunc template.FuncMap

	// Global data
	data map[string]string

	// Base template allways loaded first
	layoutTemplate string

	// Template extension
	templateExtension string
}

func New() *Template {
	templateTemplate := Template{
		template:          template.New("template"),
		templatePath:      "templates",
		templateFunc:      template.FuncMap{},
		data:              make(map[string]string),
		templateExtension: "html",
	}

	templateTemplate.initErrorTemplate()

	return &templateTemplate
}

func (t *Template) SetTemplatePath(path string) *Template {
	t.templatePath = path
	return t
}

func (t *Template) AddFunc(key string, function interface{}) *Template {
	t.templateFunc[key] = function
	return t
}

func (t *Template) SetLayout(layoutTemplate string) *Template {
	t.layoutTemplate = t.getPath(layoutTemplate)
	return t
}

func (t *Template) SetTemplateFileExt(ext string) *Template {
	t.templateExtension = ext
	return t
}

func (t *Template) SetDelimiters(left string, right string) *Template {
	t.template.Delims(left, right)
	return t
}

func (t *Template) Data(key string, value string) *Template {
	t.data[key] = value
	return t
}

func (t *Template) getPath(file string) string {
	filename := fmt.Sprintf("%s.%s", file, t.templateExtension)
	return path.Join(t.templatePath, filename)
}

func (t *Template) readFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (t *Template) mergeData(data map[string]interface{}) map[string]interface{} {
	dataMap := make(map[string]interface{})

	// clone global data
	for key, value := range t.data {
		dataMap[key] = value
	}

	// add data to map
	for key, value := range data {
		dataMap[key] = value
	}

	return dataMap
}

func (t *Template) initErrorTemplate() {
	t.template.New("error").Parse(`<html>
	<head>
		<title>{{ if .title }} {{ .title }} {{ else }} Server Error {{ end }}</title>
		<style>
			html, body {margin:0; padding:0; font-family: Verdana, Geneva, sans-serif;}
			body {background-color: #292929;}
			h1 {color: #fff; text-align: center;}
		</style>
	</head>
	<body>
		<h1>{{ if .message }} {{ .message }} {{ else }} Server Error {{ end }}</h1>
	</body>
	</html>`)
}

func (t *Template) AddTemplates(templates ...string) error {
	// load layout template
	templateString, err := t.readFile(t.layoutTemplate)
	if err != nil {
		return err
	}

	_, err = t.template.Parse(templateString)
	if err != nil {
		return err
	}

	for _, item := range templates {
		tPath := t.getPath(item)
		templateString, err := t.readFile(tPath)
		if err != nil {
			return fmt.Errorf("error reading template file: %s", tPath)
		}
		newTmpl := t.template.New(item)
		_, err = newTmpl.Parse(templateString)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) Execute(wr io.Writer, data map[string]interface{}) error {
	tempData := t.mergeData(data)
	return t.template.ExecuteTemplate(wr, "template", tempData)
}

func (t *Template) ExecuteTemplate(wr io.Writer, template string, data map[string]interface{}) error {
	tempData := t.mergeData(data)
	return t.template.ExecuteTemplate(wr, template, tempData)
}
