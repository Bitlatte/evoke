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
	Process(asset *Asset) (*Asset, error)
}
