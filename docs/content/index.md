# Evoke: The Powerful Little Static Site Generator

Welcome to Evoke, a static site generator that is small, fast, and powerful. Evoke is designed to be easy to use, with sensible defaults that allow for near-zero configuration. It also has full template support and an extension system that allows you to extend the core functionality.

## Why Evoke?

There are a lot of static site generators out there, so why build another one? The answer is simple: I wanted a static site generator that was small, fast, and powerful, but also easy to use. I wanted a tool that would let me build a website without having to worry about a lot of configuration or setup. I also wanted a tool that was extensible, so I could add new features and functionality as needed.

Evoke is the result of that vision. It's a tool that I use every day to build my own websites, and I hope that you'll find it as useful as I do.

## Core Concepts

Evoke is built around a few core concepts that make it powerful and easy to use. Understanding these concepts will help you get the most out of Evoke.

- **[Build Process](./core-concepts/build-process.md):** Learn how Evoke takes your content and turns it into a static website.
- **[Configuration](./core-concepts/configuration.md):** Discover how to customize your project with the `evoke.yaml` file.
- **[Content](./core-concepts/content.md):** Find out how to create and organize your content using HTML and Markdown.
- **[Directory Structure](./core-concepts/directory-structure.md):** Understand the purpose of each directory in an Evoke project.
- **[Partials](./core-concepts/partials.md):** Learn how to create reusable snippets of HTML.
- **[Performance](./core-concepts/performance.md):** See what makes Evoke so fast and efficient.

## Quick Start

Ready to get started? Here's a quick overview of how to install Evoke and create your first project.

### Installation

To install Evoke, you'll need to have Go installed on your system. You can then install Evoke using the following command:

```bash
go install github.com/Bitlatte/evoke/cmd/evoke@latest
```

### Create a New Project

To create a new project, simply create a new directory and add a `content` directory inside it.

```bash
mkdir my-project
cd my-project
mkdir content
```

### Create Some Content

Create a new file in the `content` directory called `index.md` and add some content to it.

```markdown
# Hello, World!

This is my first Evoke page.
```

### Build Your Site

To build your site, run the `evoke build` command from the root of your project.

```bash
evoke build
```

Evoke will build your site and place the output in the `dist` directory.

## Getting Started

If you're ready to learn more, head over to the [Getting Started](./getting-started.md) page for a more in-depth guide to creating your first project.
