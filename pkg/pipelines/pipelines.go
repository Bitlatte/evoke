// Package pipelines provides the content processing pipelines for evoke.
package pipelines

import "io"

// Asset represents a file being processed by the pipeline.
type Asset struct {
	Path     string
	Content  io.Reader
	Metadata map[string]interface{}
}

// Pipeline is an interface for processing assets.
type Pipeline interface {
	Name() string
	Process(asset *Asset) (*Asset, error)
}
