// Package pipelines provides the content processing pipelines for evoke.
package pipelines

import (
	"bytes"
	"path/filepath"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

// MarkdownPipeline is a pipeline for processing Markdown files.
type MarkdownPipeline struct {
	Goldmark goldmark.Markdown
}

// NewMarkdownPipeline creates a new MarkdownPipeline.
func NewMarkdownPipeline(gm goldmark.Markdown) *MarkdownPipeline {
	return &MarkdownPipeline{Goldmark: gm}
}

// Name returns the name of the pipeline.
func (p *MarkdownPipeline) Name() string {
	return "markdown"
}

// Process processes the asset.
func (p *MarkdownPipeline) Process(asset *Asset) (*Asset, error) {
	if filepath.Ext(asset.Path) != ".md" {
		return asset, nil
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(asset.Content)
	if err != nil {
		return nil, err
	}

	frontMatter, body, err := p.parseFrontMatter(buf.Bytes())
	if err != nil {
		return nil, err
	}

	output := new(bytes.Buffer)
	if err := p.Goldmark.Convert(body, output); err != nil {
		return nil, err
	}

	asset.Content = output
	asset.Metadata = frontMatter
	asset.Path = asset.Path[:len(asset.Path)-3] + ".html"

	return asset, nil
}

// parseFrontMatter parses the front matter from the content.
func (p *MarkdownPipeline) parseFrontMatter(content []byte) (map[string]interface{}, []byte, error) {
	var frontMatter map[string]interface{}
	body := content

	if bytes.HasPrefix(content, []byte("---")) {
		end := bytes.Index(content[3:], []byte("---"))
		if end != -1 {
			err := yaml.Unmarshal(content[3:end+3], &frontMatter)
			if err != nil {
				return nil, nil, err
			}
			body = bytes.TrimSpace(content[end+6:])
		}
	}

	return frontMatter, body, nil
}
