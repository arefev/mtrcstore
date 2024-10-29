package service

import (
	"html/template"
	"io"
)

func ListHTML(w io.Writer, list map[string]string) error {
	tpl, err := ListTpl()
	if err != nil {
		return err
	}

	data := struct {
		Title string
		Items map[string]string
	}{
		Title: "List of metrics",
		Items: list,
	}

	if err := tpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}

func ListTpl() (*template.Template, error) {
	tpl := `
	<html lang="ru">
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>{{.Title}}</title>
		</head>
		<body>
			<ul>
				{{ range $key, $value := .Items }}
				<li><strong>{{ $key }}</strong>: {{ $value }}</li>
				{{ end }}
			</ul>
		</body>
	</html>
	`
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		return t, err
	}

	return t, nil
}
