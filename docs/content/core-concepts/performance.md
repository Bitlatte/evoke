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

## Benchmark Results

To provide a transparent look at our performance, we've included the results from our benchmark tests. These tests were run on an Apple M1 CPU and measure the performance of key components in the Evoke ecosystem.

The following metrics are used:

-   `ns/op`: The average time each operation takes in nanoseconds. Lower is better.
-   `B/op`: The average number of bytes allocated per operation. Lower is better.
-   `allocs/op`: The average number of memory allocations per operation. Lower is better.

### Build (`pkg/build`)

This benchmark measures the time it takes to run the entire build process from start to finish. This includes loading plugins, copying public assets, processing content, and running all associated hooks. It provides a holistic view of the site generation time.

| Benchmark      | Iterations | Time/op (ns/op) | Memory/op (B/op) | Allocations/op |
| -------------- | ---------- | --------------- | ---------------- | -------------- |
| BenchmarkBuild | 3003       | 380991          | 100121           | 489            |

### Content (`pkg/content`)

These benchmarks measure the time it takes to process different sizes of HTML and Markdown files. The "small," "medium," and "large" variants correspond to files of approximately 1KB, 10KB, and 100KB, respectively. This helps to understand how Evoke's performance scales with content size.

| Benchmark                       | Iterations | Time/op (ns/op) | Memory/op (B/op) | Allocations/op |
| ------------------------------- | ---------- | --------------- | ---------------- | -------------- |
| BenchmarkProcessHTML/Small      | 21298      | 56975           | 848              | 22             |
| BenchmarkProcessHTML/Medium     | 22236      | 64018           | 912              | 22             |
| BenchmarkProcessHTML/Large      | 22710      | 61271           | 976              | 22             |
| BenchmarkProcessMarkdown/Small  | 22510      | 60711           | 6162             | 40             |
| BenchmarkProcessMarkdown/Medium | 19141      | 60015           | 6684             | 44             |
| BenchmarkProcessMarkdown/Large  | 19557      | 59026           | 7180             | 48             |

### Partials (`pkg/partials`)

This benchmark measures the time it takes to load and parse all partial templates from the `partials` directory. These templates are cached in memory after the first load, so this benchmark reflects the initial setup cost.

| Benchmark             | Iterations | Time/op (ns/op) | Memory/op (B/op) | Allocations/op |
| --------------------- | ---------- | --------------- | ---------------- | -------------- |
| BenchmarkLoadPartials | 29602      | 40587           | 9913             | 91             |

### Plugins (`pkg/plugins`)

This benchmark measures the overhead of the plugin system itself, without any specific plugin logic. It helps to quantify the baseline cost of having the plugin system enabled.

| Benchmark       | Iterations | Time/op (ns/op) | Memory/op (B/op) | Allocations/op |
| --------------- | ---------- | --------------- | ---------------- | -------------- |
| BenchmarkPlugin | 20985      | 63041           | 9511             | 179            |

### Util (`pkg/util`)

These benchmarks measure the performance of common file system operations, such as copying single files and entire directories. This is crucial for understanding the performance of asset handling.

| Benchmark                | Iterations | Time/op (ns/op) | Memory/op (B/op) | Allocations/op |
| ------------------------ | ---------- | --------------- | ---------------- | -------------- |
| BenchmarkCopyFile        | 17456      | 75932           | 33345            | 10             |
| BenchmarkCopyDirectory   | 6166       | 195820          | 72500            | 63             |

## Comparative Analysis

To provide a clear picture of how Evoke stacks up against other popular static site generators, we've conducted a comparative analysis with Hugo. The following benchmarks were run on a test site with 100 markdown files of approximately 1KB each.

| SSG   | Build Time (real) | Peak Memory |
| ----- | ----------------- | ----------- |
| Evoke | 0.08s             | 10.4 MB     |
| Hugo  | 0.39s             | 25.4 MB     |

As the results show, Evoke is significantly faster and uses less memory than Hugo for this test case. This is a testament to Evoke's lightweight architecture and efficient design.

## Summary

Evoke's commitment to performance means you can iterate on your site more quickly and spend less time waiting for builds. We are continuously working to make Evoke even faster and more efficient, and we encourage our community to adopt performance-conscious practices when developing plugins.
