# Partials

Partials are reusable HTML snippets that can be included in your content files. They are stored in the `partials` directory.

## Creating a Partial

To create a partial, simply create an HTML file in the `partials` directory.

### Example

`partials/header.html`:

```html
<header>
  <h1>{{ .siteName }}</h1>
</header>
```

## Using a Partial

To use a partial in your content file, use the `template` keyword.

### Example

`content/_layout.html`:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{ .siteName }}</title>
  </head>
  <body>
    {{ template "header.html" . }}
    {{ .content }}
  </body>
</html>
