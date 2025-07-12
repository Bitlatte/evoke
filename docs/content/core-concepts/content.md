# Content

Evoke supports both Markdown and HTML for creating content, giving you the flexibility to choose the best format for your needs.

## Markdown (`.md`)

Markdown is a lightweight markup language that is perfect for writing content like blog posts, articles, and documentation. Evoke uses the Goldmark library to convert your Markdown files to HTML.

### Example

```markdown
# My First Page

This is a paragraph. I can use **bold** and *italic* text.

- This is a list item.
- This is another list item.
```

## HTML (`.html`)

For more complex layouts or when you need precise control over the output, you can use standard HTML files. Any template syntax within these files will be processed by Evoke.

### Example

```html
<h1>My First Page</h1>
<p>This is a standard HTML page.</p>
```

## Routing

Evoke creates routes based on the file and directory structure within your `content` directory. For example, consider the following structure:

```
content/
├── about.md
└── blog/
    ├── post-1.md
    └── post-2.html
```

This will generate the following pages:

- `/about.html`
- `/blog/post-1.html`
- `/blog/post-2.html`

## Frontmatter

You can add metadata to your Markdown files using YAML frontmatter. This is a block of YAML at the top of the file, enclosed in triple-dashed lines (`---`).

Frontmatter allows you to define variables that can be accessed in your templates. This is useful for setting page titles, authors, dates, and other custom data.

### Example

Here's an example of a Markdown file with frontmatter:

```markdown
---
title: "My First Blog Post"
author: "Jane Doe"
date: "2024-07-08"
tags: ["tech", "golang"]
---

# My First Blog Post

This is the content of my blog post.
```

### Accessing Frontmatter in Templates

You can access these variables in your templates using the `.Page` object. For example, to display the title and author in a layout:

```html
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Page.title }}</title>
</head>
<body>
  <h1>{{ .Page.title }}</h1>
  <p>By {{ .Page.author }} on {{ .Page.date }}</p>

  <div>
    {{ .Content }}
  </div>
</body>
</html>
```

In this example, `{{ .Content }}` is a special variable that contains the rendered HTML of the Markdown content.

Frontmatter is supported for both Markdown and HTML files.

## Layouts

Evoke uses a simple layout system to help you create consistent page structures. By default, Evoke will look for a `_layout.html` file in the same directory as your content file. If it doesn't find one, it will look in the parent directory, and so on, all the way up to the `content` directory.

### Example

Consider the following directory structure:

```
content/
├── _layout.html
└── blog/
    ├── _layout.html
    └── post-1.md
```

In this example, `post-1.md` will be rendered using the `blog/_layout.html` file. If `blog/_layout.html` didn't exist, it would be rendered using `content/_layout.html`.

This allows you to create a default layout for your entire site, and then override it for specific sections of your site.
