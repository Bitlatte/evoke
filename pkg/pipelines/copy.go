package pipelines

import (
	"io"
	"os"
	"path/filepath"
)

// CopyPipeline is a pipeline for copying files.
type CopyPipeline struct{}

// NewCopyPipeline creates a new CopyPipeline.
func NewCopyPipeline() *CopyPipeline {
	return &CopyPipeline{}
}

// Process processes the asset.
func (p *CopyPipeline) Process(asset *Asset) (*Asset, error) {
	ext := filepath.Ext(asset.Path)
	if ext == ".md" || ext == ".html" {
		return asset, nil
	}

	outputPath := filepath.Join("dist", asset.Path[len("content"):])
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return nil, err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, asset.Content)
	return asset, err
}
