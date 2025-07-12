// Package pipelines provides the content processing pipelines for evoke.
package pipelines

import (
	"path/filepath"
)

// CopyPipeline is a pipeline for copying files.
type CopyPipeline struct{}

// NewCopyPipeline creates a new CopyPipeline.
func NewCopyPipeline() *CopyPipeline {
	return &CopyPipeline{}
}

// Name returns the name of the pipeline.
func (p *CopyPipeline) Name() string {
	return "copy"
}

// Process processes the asset.
func (p *CopyPipeline) Process(asset *Asset) (*Asset, error) {
	ext := filepath.Ext(asset.Path)
	if ext == ".md" || ext == ".html" {
		return asset, nil
	}

	return asset, nil
}
