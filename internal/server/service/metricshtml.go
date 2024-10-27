package service

import "fmt"

func MetricsHtml(list map[string]float64) string {
	var html string
	var li string
	for name, val := range list {
		li += fmt.Sprintf("<li>%s: %f</li>", name, val)
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
