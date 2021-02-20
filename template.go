package skm

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

// Template Type
type Template struct {

	// logging
	logg *log.Logger

	// base directory to locate all template files.
	templateDir string

	// base url
	baseURL string

	// staticUrl for loading js,css & images
	staticURL string

	// variables for us in all template files
	globalData map[string]string

	// left delimter
	left string

	// right delimiter
	right string

	// base template load first
	layoutTemplate string

	// template extsion for cleaner loading of files.
	templateExtension string
}

// NewTemplate factory
func NewTemplate(
	logg *log.Logger,
	templateDir string,
	baseURL string,
	staticURL string,
) *Template {
	return &Template{
		logg:              logg,
		templateDir:       templateDir,
		baseURL:           baseURL,
		staticURL:         staticURL,
		left:              "{{",
		right:             "}}",
		templateExtension: "html",
		globalData:        make(map[string]string),
	}
} // end New()

// SetLayout template file
func (t *Template) SetLayout(layoutTemplate string) {
	t.layoutTemplate = t.getTemplatePath(layoutTemplate)
} // end SetLayout()

// SetDelimiters change the left and right delimiters for all templates
func (t *Template) SetDelimiters(left string, right string) {
	t.left = left
	t.right = right
} //end SetDelimiters()

// SetFileExtension to be appended to the filename to locate the template
// files
func (t *Template) SetFileExtension(ext string) {
	t.templateExtension = ext
} // end SetFileExtension()

// SetData set global data to be applied to all templates
func (t *Template) SetData(key string, value string) {
	t.globalData[key] = value
} // end SetData()

// StaticURL template function
func (t *Template) StaticURL(uri string) string {
	if strings.HasPrefix(uri, "/") {
		return fmt.Sprintf("%s%s", t.staticURL, uri)
	}
	return fmt.Sprintf("%s/%s", t.staticURL, uri)
} // end StaticURL()

// BaseURL template function
func (t *Template) BaseURL(uri string) string {
	if strings.HasPrefix(uri, "/") {
		return fmt.Sprintf("%s%s", t.baseURL, uri)
	}
	return fmt.Sprintf("%s/%s", t.baseURL, uri)
} // end BaseURL()

func (t *Template) getTemplatePath(file string) string {
	filename := fmt.Sprintf("%s.%s", file, t.templateExtension)
	return path.Join(t.templateDir, filename)
} // end getTemplatePath()

func (t *Template) readFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ErrorTemplate produces a error message template
func (t *Template) ErrorTemplate() *template.Template {
	tmpl, _ := template.New("error").Parse(`<html>
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
	return tmpl
}

// Render loads the templates and executes them.
func (t *Template) Render(writer *io.Writer, templates ...string) (*template.Template, error) {

	// new template
	tmpl := template.New("layout")
	// set delimiters
	tmpl.Delims(t.left, t.right)
	// add template functions
	tmpl.Funcs(template.FuncMap{
		"baseurl":   t.BaseURL,
		"staticurl": t.StaticURL,
	})

	// read layout template
	templateString, err := t.readFile(t.layoutTemplate)
	if err != nil {
		t.logg.Print(err)
		return t.ErrorTemplate(),
			fmt.Errorf("error reading template file: %s", t.layoutTemplate)
	}

	_, err = tmpl.Parse(templateString)
	if err != nil {
		return t.ErrorTemplate(), err
	}

	for _, item := range templates {
		tPath := t.getTemplatePath(item)
		templateString, err := t.readFile(tPath)
		if err != nil {
			return t.ErrorTemplate(),
				fmt.Errorf("error reading template file: %s", tPath)
		}
		newTmpl := tmpl.New(item)
		_, err = newTmpl.Parse(templateString)
		if err != nil {
			return t.ErrorTemplate(), err
		}
	}

	return tmpl, nil
} // end Template()
