// Package pipelines provides the content processing pipelines for evoke.
package pipelines

import (
	"bytes"
	"os"
	"path/filepath"
)

// HTMLPipeline is a pipeline for processing HTML files.
type HTMLPipeline struct{}

// NewHTMLPipeline creates a new HTMLPipeline.
func NewHTMLPipeline() *HTMLPipeline {
	return &HTMLPipeline{}
}

// Name returns the name of the pipeline.
func (p *HTMLPipeline) Name() string {
	return "html"
}

// Process processes the asset.
func (p *HTMLPipeline) Process(asset *Asset) (*Asset, error) {
	if filepath.Ext(asset.Path) != ".html" {
		return asset, nil
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(asset.Content)
	if err != nil {
		return nil, err
	}

	asset.Content = buf

	outputPath := filepath.Join("dist", asset.Path[len("content"):])
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return nil, err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	_, err = outFile.Write(buf.Bytes())
	return asset, err
}
