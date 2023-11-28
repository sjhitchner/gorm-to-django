package django

import (
	// "fmt"
	"io"
	"text/template"

	"github.com/sjhitchner/gorm-to-django/gorm"
)

const Template = `
from django.db import models

{{ range . }}
class {{ .Name }}(models.Model):

	{{ range .Fields }}
	{{ .SnakeCase }} = models.{{- .Type }}(.Args)
	{{ end }}

	class Meta:
		table_name: {{ .TableName }}	

	def __str__(self):
		return self.name

{{ end }}
`

func Generate(w io.Writer, in <-chan Model) error {

	// Create a template with a unique name
	tmpl, err := template.New("django").Parse(Template)
	if err != nil {
		return err
	}

	var models []Model
	for model := range in {
		models = append(models, model)
	}

	if err := tmpl.Execute(w, models); err != nil {
		return err
	}

	return nil
}

func Convert(in <-chan gorm.Model) <-chan Model {
	out := make(chan Model)

	go func() {
		defer close(out)

		for gm := range in {
			// fmt.Println(gm)

			out <- Model{
				Name:      gm.Name,
				TableName: gm.TableName,
			}
		}

	}()
	return out
}
