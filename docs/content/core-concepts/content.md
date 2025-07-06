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

## Front Matter

You can add front matter to your content files to set variables that can be used in your layouts and partials. Front matter is written in YAML and is placed at the top of your content file, enclosed in triple-dashed lines.

### Example

```yaml
---
title: "My First Page"
---

# {{ .title }}

This is my first page.
```