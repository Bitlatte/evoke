# Development Server

Evoke comes with a powerful development server that makes it easy to preview your site locally and see changes in real-time. To start the server, run the following command in your project's root directory:

```bash
evoke serve
```

This will start a local server, typically at `http://localhost:8990`, and watch your project files for changes.

## Live Reloading

The development server features live reloading, which means that it will automatically reload your browser whenever you make a change to a file. This is a huge productivity booster, as it allows you to see the results of your changes instantly without having to manually refresh the page.

### How It Works

The development server uses a WebSocket connection to communicate with your browser. When you start the server, it injects a small JavaScript file into each HTML page. This script establishes a WebSocket connection with the server and listens for messages.

When you change a file, the server detects the change and sends a message to the browser over the WebSocket connection. The browser then reloads the page to reflect the changes.

## CSS Hot-Reloading

For an even faster development experience, the development server supports CSS hot-reloading. This means that when you change a CSS file, the new styles are injected directly into the page without a full page reload. This is especially useful when you're tweaking the design of your site, as it allows you to see the results of your changes instantly.

## Error Overlay

If you make a mistake in your code that causes the build to fail, the development server will display an error overlay in your browser. This overlay shows the error message and the file that caused the error, making it easy to identify and fix the problem.
