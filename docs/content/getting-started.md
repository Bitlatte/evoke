---
title: "Getting Started with Evoke"
---

# Getting Started with Evoke

Welcome to Evoke! This guide will walk you through installing Evoke and creating your first website.

## Installation

You can install Evoke in one of the following ways:

**1. Installer Script (Recommended):**

You can use the following command to download and install the latest version of Evoke for your system:

```bash
curl -sSL https://raw.githubusercontent.com/Bitlatte/evoke/main/install/install.sh | sh
```

**2. From a Release:**

Download the pre-compiled binary for your operating system from the [latest release](https://github.com/Bitlatte/evoke/releases/latest) on GitHub. Unzip the archive and place the `evoke` binary in a directory that is in your system's `PATH`.

**3. From Source:**

If you have Go installed, you can also install Evoke from source using the `go install` command:

```bash
go install github.com/Bitlatte/evoke/cmd/evoke@latest
```

## Your First Project

Evoke is designed to work with minimal setup. Here's how to create a basic site:

1.  **Create a Project Directory:**

    ```bash
    mkdir my-awesome-site
    cd my-awesome-site
    ```

2.  **Add a Content Directory:**

    This is the only directory you need to get started.

    ```bash
    mkdir content
    ```

3.  **Create Your First Page:**

    Create a file named `index.md` inside the `content` directory:

    ```markdown
    # Welcome to My Awesome Site!

    This is my first page. I can use **Markdown** to format my text.
    ```

4.  **Build Your Site:**

    Run the `evoke build` command from your project's root directory:

    ```bash
    evoke build
    ```

    Evoke will generate your static site in a new `dist` directory. Open `dist/index.html` in your browser to see the result.

## Project Structure Explained

As your project grows, you can add more directories to organize your files:

```
.
├── content/      # Your site's pages (Markdown or HTML)
├── partials/     # Reusable HTML snippets
├── public/       # Static assets (CSS, images, etc.)
├── extensions/   # Custom Evoke extensions
└── evoke.yaml    # Optional configuration file
```

- **`content/`**: This is where all your website's pages live. Evoke processes these files and converts them to HTML.
- **`partials/`**: This directory holds reusable HTML snippets that you can include in your pages, like headers or footers.
- **`public/`**: Any files in this directory (e.g., CSS, JavaScript, images) are copied directly to the `dist` folder without changes.
- **`extensions/`**: You can add custom Go code here to extend Evoke's functionality.
- **`evoke.yaml`**: This optional file allows you to customize your site's settings.

## What's Next?

You've successfully built your first site with Evoke! To learn more about what you can do, check out the **[Core Concepts](./core-concepts/build-process.md)**.
