# Performance

Evoke is engineered from the ground up for speed and efficiency, ensuring that your site builds quickly without monopolizing system resources. This is accomplished through a combination of a lightweight architecture, efficient algorithms, and the inherent performance of the Go programming language.

## The Go Advantage

Evoke is written in Go, a language celebrated for its performance and concurrency. This choice provides several key advantages:

-   **Single Binary Deployment:** Evoke compiles to a single, self-contained binary. This means there are no external dependencies to install or manage, making it incredibly fast to deploy and execute.
-   **Built-in Concurrency:** Go's goroutines and channels provide a powerful and efficient model for concurrency. Evoke leverages this to parallelize tasks and maximize the use of available CPU cores.

## Core Performance Features

### High-Throughput Build Process

Evoke's build process is designed for maximum throughput. Here are some of the key features that make it so fast:

-   **Parallel File Processing:** Evoke processes your content files in parallel, taking full advantage of multi-core processors to dramatically reduce build times. It creates a pool of workers, with one worker per CPU core, to ensure that your site is built as quickly as possible.
-   **In-Memory Caching:** Layouts and templates are parsed once and then cached in memory. This avoids redundant file I/O and parsing operations, resulting in a significant speed boost. The cache is implemented using a `sync.Map`, which is optimized for concurrent access.
-   **Efficient Memory Management:** Evoke is designed to be light on memory usage. We use a `sync.Pool` to reuse memory buffers for file I/O and content processing. This reduces the number of memory allocations and the pressure on the garbage collector, leading to faster and more consistent build times.
-   **Singleton Parsers:** The Goldmark Markdown parser is initialized only once and then reused for all Markdown files. This avoids the significant overhead of creating a new parser for each file.

### The Plugin System and Performance

Evoke's plugin system is designed to be flexible and powerful, but it's important to be aware of the performance implications of the plugins you use. While the core of Evoke is highly optimized, a poorly written plugin can slow down your build.

Here are some things to keep in mind when writing or using plugins:

-   **Plugin Hooks:** Plugins can "hook" into various stages of the build process. Be mindful of the hooks you use and the work you do in them. For example, a heavy computation in the `OnContentLoaded` hook will be executed for every single file, which can have a significant impact on build times.
-   **Memory Allocations:** Be mindful of memory allocations in your plugins. If you need to work with large amounts of data, consider using a `sync.Pool` to reuse buffers, just like Evoke does internally.
-   **Caching:** If your plugin performs expensive operations, consider implementing your own caching layer to avoid redundant work.

## Summary

Evoke's commitment to performance means you can iterate on your site more quickly and spend less time waiting for builds. We are continuously working to make Evoke even faster and more efficient, and we encourage our community to adopt performance-conscious practices when developing plugins.
