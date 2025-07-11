package pipelines

import (
	"bytes"

	"github.com/Bitlatte/evoke/pkg/plugins"
	"github.com/Bitlatte/evoke/proto"
)

// GRPC is a pipeline that uses a gRPC plugin to process assets.
type GRPC struct {
	Plugin plugins.Plugin
	name   string
}

// NewGRPCPipeline creates a new gRPC pipeline.
func NewGRPCPipeline(plugin plugins.Plugin, name string) *GRPC {
	return &GRPC{
		Plugin: plugin,
		name:   name,
	}
}

// Name returns the name of the pipeline.
func (p *GRPC) Name() string {
	return p.name
}

// Process processes an asset using the gRPC plugin.
func (p *GRPC) Process(asset *Asset) (*Asset, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(asset.Content); err != nil {
		return nil, err
	}

	processedAsset, err := p.Plugin.ProcessAsset(&proto.Asset{
		Path:         asset.Path,
		Content:      buf.Bytes(),
		PipelineName: p.name,
	})
	if err != nil {
		return nil, err
	}

	return &Asset{
		Path:    processedAsset.Path,
		Content: bytes.NewReader(processedAsset.Content),
	}, nil
}
