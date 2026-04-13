package templates

import (
	"bytes"
	"text/template"
)

type GoTemplateEngine struct {
	predefined map[string]string
}

func NewGoTemplateEngine() *GoTemplateEngine {
	return &GoTemplateEngine{
		predefined: map[string]string{
			"welcome": "Hi {{.name}}, welcome to our platform!",
			"alert":   "ALERT: {{.message}}. Severity: {{.level}}",
		},
	}
}

func (t *GoTemplateEngine) Render(templateName string, customBody string, data map[string]interface{}) (string, error) {
	var tmplString string

	if customBody != "" {
		tmplString = customBody
	} else {
		tmplString = t.predefined[templateName]
	}

	tmpl, err := template.New("notification").Parse(tmplString)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
