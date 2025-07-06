# evoke

A powerful little static site generator.

## Overview

The purpose of evoke is to be a small, yet powerful static site generator. This is achived through the following methods:

- Sensible defaults allowing for near zero configuration.
- Complete template support with no opinions.
- Extension system for extending the core functionality.

There are more things we could mention but I think its best to let you experience it for yourself.

## Usage

- `evoke build`: builds your content into static HTML.
- `evoke extension get [url]`: get a new extension from a url.

## Getting Started

### Installation

**NOT COMPLETED**

### Project Structure

An evoke project is simple to get started. All you need is the following folder structure and you're good to go. Here is an example of an evoke project:

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
└── templates
    ├── base.html
    └── post.html
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

### Templates Directory

This directory allows you to define templates for your content. It also allows for HTML partials.

**base.html**: This is your root layout in a sense. Whatever is defined here is going to show up on every page.
**post.html**: This is where things get interesting, notice how there was a posts directory before? this file is the template for every file under the posts directory. Say we had an `articles` directory instead, you could then use `article.html` in the templates. This allows for very flexible content creation as you can basically define whatever you want.

#### Partials

HTML partials are just snippets that get reused in multiple places for example a button. These partials can be any file that does not meet the requirements listed above. It is preferred to put partials in the `templates/partials` directory but this isn't enforced. In fact you can just mash all your partials right next to your templates and it will work just fine. It's just easier to manage by organizing them.

## evoke.yaml File

This file is optional. If you stick to all the defaults, there is no real need for this file. If you do need an evoke.yaml file, just know it accepts any key-value pair you need. You will have access to those key-values inside your templates like such:

{{ ._key_ }}

where _key_ is the key associated with a value.

Outside of this there arent many rules that need followed by the core program. Although another interesting thing evoke has is its extension system.

## Extensions

Evoke allows for custom extensions to be loaded on a per project basis. Just add a `extensions` folder to the project and add plugins using the command line. Extensions are typically loaded from a url using the cli:

```
evoke extension get [url]
```

This will automatically pull the extension and build it. Extensions can hook into the following:
- BeforeBuild: extension will run before the core build process
- AfterBuild: extension will run after the core build process

While not currently available, there are plans for an "extension library" if you will. Basically just a place with tons of extensions.

## Wraping Up

Yeah so this pretty much sums everything you need to know up. Only one last thing you need to be aware of, its "evoke" not "Evoke". Have a great day.
