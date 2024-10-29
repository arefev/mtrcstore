package service

import "fmt"

func MetricsHTML(list map[string]string) string {
	var html string
	var li string
	for name, val := range list {
		li += fmt.Sprintf("<li>%s: %s</li>", name, val)
	}

	html = fmt.Sprintf(`
	<html lang="ru">
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Список метрик</title>
		</head>
		<body>
			<ul>%s</ul>
		</body>
	</html>
	`, li)

	return html
}
