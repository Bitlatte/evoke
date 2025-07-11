# Plugin Service Definition

This document provides the full definition of the `Plugin` gRPC service, which is used to create plugins for Evoke.

## Service Definition

```protobuf
syntax = "proto3";

package proto;

option go_package = "github.com/Bitlatte/evoke/proto";

// The main service that plugins must implement.
service Plugin {
  // --- General Build Hooks ---

  // Called once before the entire build process begins.
  // Useful for setup tasks or pre-build validation.
  rpc OnPreBuild(OnPreBuildRequest) returns (OnPreBuildResponse);

  // Called after the configuration file (evoke.yaml) is loaded.
  // Allows plugins to read or even modify the configuration.
  rpc OnConfigLoaded(OnConfigLoadedRequest) returns (OnConfigLoadedResponse);

  // Called after the 'public' directory has been copied to 'dist'.
  rpc OnPublicAssetsCopied(OnPublicAssetsCopiedRequest) returns (OnPublicAssetsCopiedResponse);

  // --- Content Processing Hooks ---

  // Called for each content file after it's read from disk but before any processing.
  // Allows modification of the raw file content.
  rpc OnContentLoaded(OnContentLoadedRequest) returns (OnContentLoadedResponse);

  // Called before the Markdown (or other format) content is rendered to HTML.
  // A plugin could use this to implement a custom renderer.
  rpc OnContentRender(OnContentRenderRequest) returns (OnContentRenderResponse);

  // Called after content is rendered to HTML but before it's placed in a layout.
  // Useful for post-processing the core HTML content.
  rpc OnHTMLRendered(OnHTMLRenderedRequest) returns (OnHTMLRenderedResponse);

  // --- Finalization Hooks ---

  // Called once after all content has been processed and written to disk.
  rpc OnPostBuild(OnPostBuildRequest) returns (OnPostBuildResponse);
}

// Represents a file being processed. This message will be reused for
// multiple hooks to pass content back and forth.
message ContentFile {
  // The relative path of the file from the content directory.
  string path = 1;
  // The raw or processed content of the file.
  bytes content = 2;
}

// Represents an asset being processed by a custom pipeline.
message Asset {
  string path = 1;
  bytes content = 2;
  string pipeline_name = 3;
}

// Represents a custom pipeline that can be registered by a plugin.
message Pipeline {
  string name = 1;
  repeated string extensions = 2;
}

message OnPreBuildRequest {}
message OnPreBuildResponse {}

message OnConfigLoadedRequest {
  bytes config = 1;
}
message OnConfigLoadedResponse {
  bytes config = 1;
}

message OnPublicAssetsCopiedRequest {}
message OnPublicAssetsCopiedResponse {}

message OnContentLoadedRequest {
  string path = 1;
  bytes content = 2;
}
message OnContentLoadedResponse {
  bytes content = 1;
}

message OnContentRenderRequest {
  string path = 1;
  bytes content = 2;
}
message OnContentRenderResponse {
  bytes content = 1;
}

message OnHTMLRenderedRequest {
  string path = 1;
  bytes content = 2;
}
message OnHTMLRenderedResponse {
  bytes content = 1;
}

message OnPostBuildRequest {}
message OnPostBuildResponse {}
