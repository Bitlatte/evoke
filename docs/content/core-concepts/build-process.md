# Build Process

When you run `evoke build`, the following steps are performed:

1.  **Load Configuration:** Evoke loads the configuration from the `evoke.yaml` file. This file contains all of the settings for your project, such as the name of your site, the URL of your site, and any custom data that you want to make available to your templates.

2.  **Create Output Directory:** Evoke creates the `dist` directory if it doesn't already exist. This is where your static site will be generated.

3.  **Copy Public Directory:** Evoke copies the contents of the `public` directory to the `dist` directory. This is where you should put any static assets that you want to be copied to your site, such as images, CSS files, and JavaScript files.

4.  **Load Partials:** Evoke loads any partials from the `partials` directory. Partials are small snippets of HTML that can be reused across multiple pages. For example, you might have a partial for your site's header and another for your site's footer.

5.  **Load Plugins:** Evoke loads any plugins from the `plugins` directory. Plugins are small programs that can be used to extend the functionality of Evoke. For example, you could use a plugin to add support for a new templating language or to add a custom build step.

6.  **Run BeforeBuild Hooks:** Evoke runs the `BeforeBuild` hook for each loaded plugin. This allows plugins to perform any necessary setup before the build process begins.

7.  **Process Content:** Evoke processes all of the content in the `content` directory. This is where you should put all of the pages for your site. Evoke supports both Markdown and HTML files.

8.  **Run AfterBuild Hooks:** Evoke runs the `AfterBuild` hook for each loaded plugin. This allows plugins to perform any necessary cleanup after the build process is complete.

## Incremental Builds

To improve build times, Evoke uses an incremental build process. This means that it only rebuilds files that have changed since the last build. This is accomplished by storing a cache of file hashes in the `dist/.cache` file.

When you run `evoke build`, Evoke first builds a dependency graph of all the files in your `content` and `partials` directories. It then compares the hashes of the files in the dependency graph to the hashes in the cache. If a file's hash has changed, or if the file is not in the cache, Evoke will rebuild the file and any files that depend on it.

This process is completely automatic and requires no configuration. However, if you ever need to force a full rebuild, you can do so by deleting the `dist/.cache` file or by running the build with the `--clean` flag.
