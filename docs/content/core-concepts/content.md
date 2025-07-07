# Content

Evoke supports both Markdown and HTML for creating content.

## Markdown

You can create content using Markdown files with the `.md` extension. Evoke uses the Goldmark library to convert Markdown to HTML.

### Example

```markdown
# My First Page

This is my first page.
```

## HTML

You can also create content using HTML files with the `.html` extension.

### Example

```html
<h1>My First Page</h1>

<p>This is my first page.</p>
```

## Frontmatter

You can add frontmatter to your Markdown files to set variables that can be used in your layouts and partials. Frontmatter is written in YAML and is placed at the top of your content file, enclosed in triple-dashed lines (`---`).

All variables defined in the frontmatter are added to the page's context and can be accessed in your templates. For example, if you define a `title` in your frontmatter, you can use `{{ .title }}` in your layout to display it.

### Example

Here's an example of how to use frontmatter to set a title and a custom variable `author`:

```yaml
---
title: "My First Page"
author: "John Doe"
---
```

These variables can then be accessed from your `_layout.html` like so:

```html
...
<title>{{ .title }}</title>
...
<body>
  ...
  <span>{{ .author }}</span>
</body>
```

In this example, {{ .title }} will be replaced with "My First Page" and {{ .author }} will be replaced with "John Doe".

Frontmatter is only supported for Markdown files and is ignored in HTML files.
