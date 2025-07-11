# evoke

A powerful little static site generator.

## Overview

The purpose of evoke is to be a small, yet powerful static site generator. This is achived through the following methods:

- Sensible defaults allowing for near zero configuration.
- Complete template support with no opinions.
- Plugin system for extending the core functionality.

There are more things we could mention but I think its best to let you experience it for yourself.

## Usage

- `evoke build`: builds your content into static HTML.
- `evoke serve`: builds and serves the site on a local development server.
- `evoke init`: initializes a new Evoke project.

## Getting Started

### Installation

You can install Evoke in one of the following ways:

**1. Installer Script (Recommended):**

You can use the following command to download and install the latest version of Evoke for your system:

```bash
curl -sSL https://raw.githubusercontent.com/Bitlatte/evoke/main/install.sh | sh
```

**2. From a Release:**

Download the pre-compiled binary for your operating system from the [latest release](https://github.com/Bitlatte/evoke/releases/latest) on GitHub. Unzip the archive and place the `evoke` binary in a directory that is in your system's `PATH`.

**3. From Source:**

If you have Go installed, you can also install Evoke from source using the `go install` command:

```bash
go install github.com/Bitlatte/evoke/cmd/evoke@latest
```

### Project Structure

The easiest way to start a new Evoke project is to use the `init` command.

```bash
mkdir my-project
cd my-project
evoke init
```

This will create a basic project structure for you. Here is an example of an evoke project:

```
.
├── content
│   ├── about.html
│   └── posts
│       ├── _index.html
│       ├── post-1.md
│       └── post-2.md
├── evoke.yaml
├── public
│   ├── css
│   │   └── style.css
│   ├── img
│   │   └── sample.jpg
│   └── js
│       └── script.js
└── partials
    ├── header.html
    └── footer.html
```

### Content Directory

The content directory holds all your content. It has a few rules but for the most part anything you put in here will be included in the final build.

- Rule 1: The content directory will define routes based on naming and folder structure. Take for example:

```
.
├── content
    ├── about.html
    └── posts
        ├── index.html
        ├── post-1.md
        └── post-2.md
```

This will create the following routes:

- about.html
- posts/index.html
- posts/post-1.html
- posts/post-2.html

Notice how we have mixed HTML and markdown in the content directories? This is to allow more advanced sites to be made. HTML files will have any template strings expanded and the rest of the file will just be the same. Markdown files on the other hand will utilize a special template which we will talk about later.

### Public Directory

This directory will simply be copied to the dist folder when building. This is so you can include images, css, javascript, or whatever in your pages.

### Partials Directory

This directory allows you to define reusable HTML snippets, called partials, for your content.

Partials are included in your layout and content files using the `template` keyword. For example, to include a partial named `header.html`, you would use the following syntax:

```html
{{ `{{ template "header.html" . }}` }}
```

## evoke.yaml File

This file is optional. If you stick to all the defaults, there is no real need for this file. If you do need an evoke.yaml file, just know it accepts any key-value pair you need. You will have access to those key-values inside your templates like such:

{{ `{{ .key }}` }}

where _key_ is the key associated with a value.

Outside of this there arent many rules that need followed by the core program. Although another interesting thing evoke has is its extension system.

## Plugins

Evoke allows for custom plugins to be loaded on a per project basis. Just add a `plugins` folder to the project and add your plugins.

Plugins can hook into the following:
- BeforeBuild: plugin will run before the core build process
- AfterBuild: plugin will run after the core build process

## Wraping Up

Yeah so this pretty much sums everything you need to know up. Only one last thing you need to be aware of, its "evoke" not "Evoke". Have a great day.
