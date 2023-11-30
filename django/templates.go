package django

const ModelTemplate = `
from django.db import models

{{ range . }}
class {{ .Name }}(models.Model):

	{{- range .Fields }}
	{{ .Name }} = models.{{- .Type }}()
	{{- end }}

	class Meta:
        managed = False
        db_table = {{ .TableName }}

	def __str__(self):
		return self.name

{{ end }}
`

const AdminTemplate = `
from django.contrib import admin

from .models import (
{{ range . }}
	{{ .Name }}Event,
{{ end }}
)

{{ range . }}
class {{ .Name }}Admin(admin.ModelAdmin):
    list_display = ['name']
	readonly_fields = 
{{ end }}

{{ range . }}
admin.site.register({{ .Name }}, {{ .Name }}Admin)
{{ end }}
`
