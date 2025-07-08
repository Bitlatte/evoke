# Performance

Evoke is designed to be fast and efficient, allowing you to build your site quickly without consuming excessive system resources. This is achieved through a combination of a lightweight design, efficient algorithms, and the power of the Go programming language.

## Why Go?

Evoke is written in Go, a language renowned for its performance and concurrency features. This choice allows Evoke to:

- **Compile to a single binary:** This means no external dependencies are needed to run Evoke, making it fast to install and execute.
- **Excellent concurrency support:** Go's built-in goroutines and channels make it easy to write highly concurrent code, which is key to Evoke's parallel processing capabilities.

## Key Performance Features

### Build Speed

Evoke's build process is highly optimized for speed. On a modern machine, a typical site builds in a fraction of a second. This is achieved through several techniques:

-   **Parallel Processing:** Evoke processes content files in parallel, taking full advantage of multi-core processors to speed up the build.
-   **Efficient Caching:** Layouts and templates are cached in memory after they are first loaded, avoiding redundant file I/O and parsing.

### Memory Usage

We have put significant effort into optimizing Evoke's memory usage to ensure it runs smoothly even on systems with limited RAM. Key memory optimizations include:

-   **Buffer Re-use with `sync.Pool`:** Evoke uses `sync.Pool` to reuse memory buffers for file copying and content processing. This dramatically reduces the number of memory allocations and the pressure on the garbage collector, which in turn leads to faster build times.
-   **Singleton Parsers:** The Markdown parser (Goldmark) is initialized only once and then reused throughout the build process. This avoids the significant memory overhead of creating new parsers for each file.

### Benchmark Results

Our internal benchmarks show that these optimizations have resulted in a **~30% reduction in total memory allocations** during a typical build, with a corresponding **~18% improvement in build speed**.

## Summary

Evoke's focus on performance means you can spend less time waiting for your site to build and more time creating content. We are continuously monitoring and improving Evoke's performance to ensure it remains one of the fastest and most efficient static site generators available.
