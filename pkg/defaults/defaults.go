// Package defaults provides default values for evoke.
package defaults

// Layout is the default layout.html content.
var Layout = `<!DOCTYPE html>
<html>
<head>
	<link rel="icon" href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>✨</text></svg>">
	<title>{{ .Site.Name }}</title>
</head>
<body>
	{{ .Content }}
</body>
</html>`
