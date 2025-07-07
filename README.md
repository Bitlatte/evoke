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

To get started, you'll need to have Go installed on your system. You can then install Evoke using the following command:

```bash
go install github.com/Bitlatte/evoke/cmd/evoke@latest
```

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

## Extensions

Evoke allows for custom extensions to be loaded on a per project basis. Just add a `extensions` folder to the project and add plugins using the command line. Extensions are typically loaded from a url using the cli:

```
evoke extension get [url]
```

This will automatically pull the extension and build it. Extensions can hook into the following:
- BeforeBuild: extension will run before the core build process
- AfterBuild: extension will run after the core build process

While not currently available, there are plans for an "extension library" if you will. Basically just a place with tons of extensions.

## Performance

The following benchmarks are run on a machine with an Apple M1 CPU. The benchmarks cover a range of site sizes to provide a comprehensive view of the engine's performance. Each test is run against both Evoke and Hugo for a direct comparison. The content for each page is dynamically generated to be roughly the specified number of lines.

- **Tiny:** 1 page with ~10 lines of content.
- **Small:** 100 pages, each with ~50 lines of content.
- **Medium:** 1000 pages, each with ~100 lines of content.
- **Large:** 1000 pages, each with ~500 lines of content.
- **Huge:** 10,000 pages, each with ~500 lines of content.

<!-- BENCHMARKS_START -->
<table>
<thead>
<tr>
<th>Benchmark</th>
<th>Evoke Time/op</th>
<th>Hugo Time/op</th>
<th>Evoke Memory/op</th>
<th>Hugo Memory/op</th>
<th>Evoke Allocs/op</th>
<th>Hugo Allocs/op</th>
</tr>
</thead>
<tbody>
<tr>
<td>Tiny-8</td>
<td>130.16 µs</td>
<td>81.04 ms</td>
<td>38.05 KB</td>
<td>9.74 KB</td>
<td>264</td>
<td>47</td>
</tr>
<tr>
<td>Small-8</td>
<td>8.11 ms</td>
<td>109.28 ms</td>
<td>3.33 MB</td>
<td>9.74 KB</td>
<td>13864</td>
<td>47</td>
</tr>
<tr>
<td>Medium-8</td>
<td>114.73 ms</td>
<td>487.20 ms</td>
<td>59.90 MB</td>
<td>9.77 KB</td>
<td>189228</td>
<td>47</td>
</tr>
<tr>
<td>Large-8</td>
<td>269.79 ms</td>
<td>2009.74 ms</td>
<td>296.96 MB</td>
<td>9.83 KB</td>
<td>596950</td>
<td>48</td>
</tr>
<tr>
<td>Huge-8</td>
<td>3219.10 ms</td>
<td>28406.37 ms</td>
<td>2945.19 MB</td>
<td>9.83 KB</td>
<td>5946592</td>
<td>48</td>
</tr>
</tbody>
</table>
<!-- BENCHMARKS_END -->

## Wraping Up

Yeah so this pretty much sums everything you need to know up. Only one last thing you need to be aware of, its "evoke" not "Evoke". Have a great day.
