{{ define "databend_create_table" -}}

USE `{{.Database}}`;

CREATE TABLE `{{.Table}}` (
	{{- range $i, $col := .Fields}}
	{{- if eq (add $i 1) (len $.Fields) }}
	{{$col.GetName}} {{$col.GetType}} {{$col.GetNull}} {{$col.GetDefault}}
	{{- else }}
	{{$col.GetName}} {{$col.GetType}} {{$col.GetNull}} {{$col.GetDefault}},
	{{- end }}
	{{- end }}
);
{{- end}}
