// Package defaults provides default values for evoke.
package defaults

// Layout is the default layout.html content.
var Layout = `<!DOCTYPE html>
<html>
<head>
	<title>{{ .Site.Name }}</title>
</head>
<body>
	{{ .Content }}
</body>
</html>`
