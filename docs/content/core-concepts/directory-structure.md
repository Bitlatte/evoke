# Directory Structure

Evoke uses a simple directory structure to organize your site.

*   `content/`: This directory contains all of your site's content, including Markdown and HTML files. It can also contain a `_layout.html` file to define the base layout for your pages. The directory structure within the `content` directory will be used to generate the URLs for your site. For example, a file at `content/blog/my-post.md` will be available at `/blog/my-post.html`.


*   `public/`: This directory contains all of your site's static assets, such as images, CSS, and JavaScript files. The contents of this directory will be copied to the `dist` directory when you build your site.

*   `partials/`: This directory contains all of your site's partials, which are reusable HTML snippets that can be included in your content files. For example, you could create a partial for your site's header and another for your site's footer.

*   `plugins/`: This directory contains all of your site's plugins, which are Go plugins that can be used to extend Evoke's functionality. For example, you could create a plugin to add support for a new templating language or to add a custom build step.

*   `dist/`: This directory is where your static site will be generated. You should not edit the contents of this directory directly, as it will be overwritten every time you build your site.

*   `evoke.yaml`: An optional configuration file for your site. This file can be used to configure your site's name, URL, and other settings.
