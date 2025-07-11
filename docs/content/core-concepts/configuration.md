# Configuration

Evoke is designed to work with zero configuration, but you can customize your project by creating an `evoke.yaml` file in the root directory.

## The `evoke.yaml` File

This file is entirely optional. If you're happy with Evoke's default settings, you don't need it. However, if you want to customize your site, this is the place to do it.

The `evoke.yaml` file uses the YAML format, which is easy to read and write. You can add any key-value pairs you need, and they will be available in your templates.

### Example

Here's an example of a more complex `evoke.yaml` file:

```yaml
siteName: "My Awesome Site"
author: "John Doe"
baseURL: "https://example.com"
social:
  twitter: "@johndoe"
  github: "johndoe"
```

### Accessing Configuration Values in Templates

All values from your `evoke.yaml` file are available in your templates under the `.Site` object. For example, to display the site name and author from the example above, you would use the following in your HTML files:

```html
<h1>{{ .Site.siteName }}</h1>
<p>By {{ .Site.author }}</p>
```

To access nested values, like the social media links, you can chain the keys:

```html
<a href="https://twitter.com/{{ .Site.social.twitter }}">Twitter</a>
<a href="https://github.com/{{ .Site.social.github }}">GitHub</a>
```

This flexibility allows you to create highly customized and dynamic templates with ease.

### Using the `Global` Key

For site-wide variables, it's a good practice to nest them under a `Global` key in your `evoke.yaml` file. This helps to keep your configuration organized.

```yaml
Global:
  siteName: "My Awesome Site"
  author: "John Doe"
```

When you use the `Global` key, you'll need to access these variables through the `.Site` object in your templates. Here’s how you would access the `siteName`:

```html
<h1>{{ .Site.Global.siteName }}</h1>
```

This approach ensures that your global variables are neatly namespaced and easy to manage.
