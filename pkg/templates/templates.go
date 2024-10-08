package templates

import (
	"bytes"
	"text/template"
)

// CompileTemplateFromMap replaces placeholders in the HTML template with values from the configMap
func CompileTemplateFromMap(tmplt string, configMap interface{}) (string, error) {
	out := new(bytes.Buffer)
	t := template.Must(template.New("compiled_template").Parse(tmplt))
	if err := t.Execute(out, configMap); err != nil {
		return "", err
	}
	return out.String(), nil
}
