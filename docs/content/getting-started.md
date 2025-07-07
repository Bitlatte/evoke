---
title: "Getting Started with Evoke"
---

# Getting Started with Evoke

Welcome to Evoke! This guide will help you get up and running with your new static site generator.

## Installation

To get started, you'll need to have Go installed on your system. You can then install Evoke using the following command:

```bash
go install github.com/Bitlatte/evoke/cmd/evoke@latest
```

## Creating Your First Site

1.  **Create a new directory for your site:**

    ```bash
    mkdir my-site
    cd my-site
    ```

2.  **Create the required directories:**

    ```bash
    mkdir content public partials extensions
    ```

3.  **Create a configuration file:**

    Create a file named `evoke.yaml` in the root of your project with the following content:

    ```yaml
    siteName: "My Awesome Site"
    ```

4.  **Create a layout file:**

    Create a file named `_layout.html` in the `content` directory with the following content:

    ```html
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>{{ .siteName }}</title>
      </head>
      <body>
        {{ .content }}
      </body>
    </html>
    ```

5.  **Create your first page:**

    Create a file named `index.md` in the `content` directory with the following content:

    ```markdown
    # Welcome to My Awesome Site!

    This is my first page.
    ```

6.  **Build your site:**

    Run the following command to build your site:

    ```bash
    evoke build
    ```

    This will generate your static site in the `dist` directory. You can now open the `dist/index.html` file in your browser to see your new site.
