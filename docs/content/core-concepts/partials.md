# Partials

Partials are reusable HTML snippets that help you keep your code DRY (Don't Repeat Yourself). They are stored in the `partials` directory and can be included in your layouts and content files.

## Creating a Partial

To create a partial, simply create an HTML file in the `partials` directory.

### Example: `partials/header.html`

```html
<header>
  <h1>{{ .Global.siteName }}</h1>
  <p>Welcome to my awesome site!</p>
</header>
```

## Using a Partial

To include a partial, use the `template` keyword. The `.` (dot) passes the current context (e.g., page variables, site configuration) to the partial.

### Example: `content/_layout.html`

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>{{ .Page.title }}</title>
</head>
<body>
  {{ template "header.html" . }}

  <main>
    {{ .Content }}
  </main>

  {{ template "footer.html" . }}
</body>
</html>
```

## Passing Custom Data to Partials

You can also pass custom data to a partial. This is useful for creating reusable components that can be customized on a per-page basis.

### Example: A `card` Partial

Let's create a partial to display a card with a title and content.

`partials/card.html`:

```html
<div class="card">
  <h2>{{ .Title }}</h2>
  <p>{{ .Content }}</p>
</div>
```

Now, you can use this partial in your content files and pass data to it using the `dict` function:

`content/index.md`:

```html
---
title: "Home Page"
---

# Welcome to the Home Page

Here are some featured items:

{{ template "card.html" (dict "Title" "Card 1" "Content" "This is the first card.") }}
{{ template "card.html" (dict "Title" "Card 2" "Content" "This is the second card.") }}
```

## Looping with Partials

Partials are also great for rendering lists of items. For example, you could loop through a list of blog posts and render a partial for each one.

### Example: Listing Blog Posts

Imagine you have a list of posts in your `evoke.yaml`:

```yaml
posts:
  - title: "Post 1"
    url: "/blog/post-1"
  - title: "Post 2"
    url: "/blog/post-2"
```

You can then loop through these posts in your template and render a partial for each one:

`content/blog.html`:

```html
<h1>Blog</h1>
<ul>
  {{ range .Global.posts }}
    {{ template "post-summary.html" . }}
  {{ end }}
</ul>
```

`partials/post-summary.html`:

```html
<li>
  <a href="{{ .url }}">{{ .title }}</a>
</li>
```

This powerful combination of partials, data, and loops allows you to build complex and maintainable websites with ease.

## Nested Partials

You can also nest partials within other partials. This is useful for creating complex components from smaller, more manageable pieces.

### Example: A `profile` Partial

Let's create a `profile` partial that uses a `card` partial.

`partials/profile.html`:

```html
<div class="profile">
  {{ template "card.html" (dict "Title" .Page.Name "Content" .Page.Bio) }}
</div>
```

Now, you can use the `profile` partial in your content files:

`content/about.md`:

```markdown
---
title: "About Me"
Name: "John Doe"
Bio: "I am a web developer."
---

# About Me

{{ template "profile.html" . }}
