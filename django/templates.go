package django

const ModelTemplate = `
from django.db import models

{{ range . }}
class {{ .Name }}(models.Model):

	{{- range .Fields }}
	{{ .Name }} = models.{{- .DjangoField }}({{ .DjangoArgs }})
	{{- end }}

	class Meta:
		managed = False
		{{- range $k, $v := .Metadata }}
		{{ $k }} = {{ $v }}
		{{- end }}

	def __str__(self):
		return self.name

{{ end }}
`

const AdminTemplate = `
from django.contrib import admin

from .models import (
{{- range . }}
	{{ .Name }},
{{- end }}
)

{{ range . }}
class {{ .Name }}Admin(admin.ModelAdmin):
	list_display = [{{ .DisplayList }}]
	readonly_fields = [{{ .ReadOnlyFields }}]
{{ end }}

{{ range . }}
admin.site.register({{ .Name }}, {{ .Name }}Admin)
{{- end }}
`
