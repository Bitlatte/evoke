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

-   **Plugin Hooks:** Plugins can "hook" into various stages of the build process. Be mindful of the hooks you use and the work you do in them. For example, a heavy computation in a hook that runs for every file can have a significant impact on build times.
-   **Memory Allocations:** Be mindful of memory allocations in your plugins. If you need to work with large amounts of data, consider using a `sync.Pool` to reuse buffers, just like Evoke does internally.
-   **Caching:** If your plugin performs expensive operations, consider implementing your own caching layer to avoid redundant work.

## Benchmark Results

To provide a transparent look at our performance, we've included the results from our benchmark tests. These tests were run on an Apple M1 CPU and measure the performance of key components in the Evoke ecosystem.

The following metrics are used:

-   `ns/op`: The average time each operation takes in nanoseconds. Lower is better.
-   `B/op`: The average number of bytes allocated per operation. Lower is better.
-   `allocs/op`: The average number of memory allocations per operation. Lower is better.

### Build (`pkg/build`)

This benchmark measures the time it takes to build a site with 100 pages from scratch. This includes loading plugins, copying public assets, processing content, and running all associated hooks. It provides a holistic view of the site generation time.

> **Note:** The clean build is now slower than it was previously. This is because Evoke now builds a dependency graph and hashes all of your files to enable incremental builds. While the initial build is slower, subsequent builds will be significantly faster.

| Benchmark      | Time/op (ms) | Memory/op (MB) | Allocations/op |
| -------------- | ------------ | -------------- | -------------- |
| BenchmarkBuild | 46.15        | 7.62           | 31283          |

### Pipelines (`pkg/pipelines`)

These benchmarks measure the time it takes for each content pipeline to process a realistic piece of content.

| Benchmark                 | Time/op (ms) | Memory/op (KB) | Allocations/op | Notes                               |
| ------------------------- | ------------ | -------------- | -------------- | ----------------------------------- |
| BenchmarkMarkdownPipeline | 0.15         | 111.13         | 429            | Processes a 100-paragraph MD file   |
| BenchmarkHTMLPipeline     | 0.01         | 32.34          | 8              | Processes a 100-paragraph HTML file |
| BenchmarkCopyPipeline     | 0.12         | 1048.63        | 2              | Processes a 1MB file                |

### Partials (`pkg/partials`)

This benchmark measures the time it takes to load and parse 50 partial templates from the `partials` directory.

| Benchmark             | Time/op (ms) | Memory/op (KB) | Allocations/op |
| --------------------- | ------------ | -------------- | -------------- |
| BenchmarkLoadPartials | 2.24         | 246.05         | 2148           |

### Plugins (`pkg/plugins`)

This benchmark measures the overhead of the plugin system by sending a 10KB payload over gRPC.

| Benchmark       | Time/op (ms) | Memory/op (KB) | Allocations/op |
| --------------- | ------------ | -------------- | -------------- |
| BenchmarkPlugin | 0.10         | 76.10          | 187            |

### Util (`pkg/util`)

These benchmarks measure the performance of common file system operations.

| Benchmark              | Time/op (ms) | Memory/op (MB) | Allocations/op | Notes                           |
| ---------------------- | ------------ | -------------- | -------------- | ------------------------------- |
| BenchmarkCopyFile      | 0.96         | 0.03           | 10             | Copies a 1MB file               |
| BenchmarkCopyDirectory | 14.06        | 2.43           | 1603           | Copies a directory with 100+ files |

## Comparative Analysis

To provide a clear picture of how Evoke stacks up against other popular static site generators, we conducted a comparative analysis with Hugo, Eleventy, and Gatsby. The following benchmarks were run on a test site with 5,000 markdown files.

The test was conducted on an Apple M1 CPU. Each project was set up with a basic configuration, and the build time was measured using the `time` command.

| SSG      | Build Time (real) | Time Difference | Times Slower |
| -------- | ----------------- | --------------- | ------------ |
| Evoke    | 1.50s             | -               | -            |
| Hugo     | 4.453s            | +2.953s         | 2.97x        |
| Eleventy | 4.650s            | +3.15s          | 3.10x        |
| Gatsby   | 22.432s           | +20.932s        | 14.95x       |

As the results show, Evoke is significantly faster than the other static site generators in this test case. This is a testament to Evoke's lightweight architecture and efficient design. While this benchmark is not exhaustive, it provides a strong indication of Evoke's performance advantages for content-heavy sites.

## Summary

Evoke's commitment to performance means you can iterate on your site more quickly and spend less time waiting for builds. We are continuously working to make Evoke even faster and more efficient, and we encourage our community to adopt performance-conscious practices when developing plugins.
