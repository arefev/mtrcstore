package service

import (
	"fmt"
	"html/template"
	"io"
)

func ListHTML(w io.Writer, list map[string]string) error {
	tpl, err := ListTpl()
	if err != nil {
		return err
	}

	data := struct {
		Items map[string]string
		Title string
	}{
		Items: list,
		Title: "List of metrics",
	}

	if err := tpl.Execute(w, data); err != nil {
		return fmt.Errorf("ListHtml template execute failed: %w", err)
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
		return t, fmt.Errorf("ListTpl template parse failed: %w", err)
	}

	return t, nil
}
