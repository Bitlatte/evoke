# Layouts

Evoke uses a simple yet powerful layout system to help you create consistent page structures for your site. Layouts are defined using `_layout.html` files, and they allow you to define a common structure for a set of pages.

## The `_layout.html` File

A layout is an HTML file that contains the basic structure of a page. It typically includes the `<html>`, `<head>`, and `<body>` tags, as well as any other common elements that you want to appear on every page, such as a header, footer, or navigation bar.

The key to a layout file is the `{{ .Content }}` variable. This is where the content of the individual pages will be injected.

### Example

Here's an example of a basic layout file:

`content/_layout.html`:
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>{{ .Page.title }}</title>
</head>
<body>
  <header>
    <h1>{{ .Global.siteName }}</h1>
  </header>

  <main>
    {{ .Content }}
  </main>

  <footer>
    <p>&copy; 2024 {{ .Global.siteName }}</p>
  </footer>
</body>
</html>
```

In this example, `{{ .Page.title }}` will be replaced with the title from the page's front matter, and `{{ .Global.siteName }}` will be replaced with the site name from the `evoke.yaml` file.

## Hierarchical Layouts

Evoke's layout system is hierarchical. When rendering a page, Evoke will look for a `_layout.html` file in the same directory as the page. If it doesn't find one, it will look in the parent directory, and so on, all the way up to the `content` directory.

This allows you to create a default layout for your entire site, and then override it for specific sections.

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

## Nested Layouts

Layouts can also be nested. This is useful for creating complex page structures with multiple levels of inheritance.

### Example

Let's say you have a base layout for your entire site, and then a separate layout for your blog that adds a sidebar.

`content/_layout.html`:
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <title>{{ .Page.title }}</title>
</head>
<body>
  {{ .Content }}
</body>
</html>
```

`content/blog/_layout.html`:
```html
<div class="container">
  <main class="main-content">
    {{ .Content }}
  </main>
  <aside class="sidebar">
    <h2>Recent Posts</h2>
    <ul>
      <li>Post 1</li>
      <li>Post 2</li>
    </ul>
  </aside>
</div>
```

When `post-1.md` is rendered, its content will first be injected into `content/blog/_layout.html` in place of `{{ .Content }}`. Then, the *entire result* of that will be injected into `content/_layout.html` in place of its `{{ .Content }}`.
