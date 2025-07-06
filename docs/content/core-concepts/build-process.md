# Build Process

When you run `evoke build`, the following steps are performed:

1.  **Load Extensions:** Evoke loads any extensions from the `extensions` directory.
2.  **Run BeforeBuild Hooks:** Evoke runs the `BeforeBuild` hook for each loaded extension.
3.  **Create Output Directory:** Evoke creates the `dist` directory if it doesn't already exist.
4.  **Copy Public Directory:** Evoke copies the contents of the `public` directory to the `dist` directory.
5.  **Load Configuration:** Evoke loads the configuration from the `evoke.yaml` file.
6.  **Load Partials:** Evoke loads any partials from the `partials` directory.
7.  **Process Content:** Evoke processes all of the content in the `content` directory.
8.  **Run AfterBuild Hooks:** Evoke runs the `AfterBuild` hook for each loaded extension.
