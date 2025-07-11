// Package plugins provides the gRPC server and client implementations for evoke plugins.
package plugins

import (
	"context"

	"github.com/Bitlatte/evoke/proto"
)

// GRPCClient is an implementation of Plugin that talks over RPC.
type GRPCClient struct {
	// Client is the gRPC client for the plugin.
	Client proto.PluginClient
}

// OnPreBuild is called before the build process starts.
func (m *GRPCClient) OnPreBuild() error {
	_, err := m.Client.OnPreBuild(context.Background(), &proto.PreBuildRequest{})
	return err
}

// OnConfigLoaded is called after the configuration is loaded.
func (m *GRPCClient) OnConfigLoaded(config []byte) ([]byte, error) {
	resp, err := m.Client.OnConfigLoaded(context.Background(), &proto.ConfigLoadedRequest{ConfigJson: string(config)})
	if err != nil {
		return nil, err
	}
	return []byte(resp.ConfigJson), nil
}

// OnPublicAssetsCopied is called after the public assets are copied.
func (m *GRPCClient) OnPublicAssetsCopied() error {
	_, err := m.Client.OnPublicAssetsCopied(context.Background(), &proto.PublicAssetsCopiedRequest{})
	return err
}

// OnContentLoaded is called after a content file is loaded.
func (m *GRPCClient) OnContentLoaded(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnContentLoaded(context.Background(), &proto.ContentFile{Path: path, Content: content})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

// OnContentRender is called after a content file is rendered.
func (m *GRPCClient) OnContentRender(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnContentRender(context.Background(), &proto.ContentFile{Path: path, Content: content})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

// OnHTMLRendered is called after the HTML is rendered.
func (m *GRPCClient) OnHTMLRendered(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnHTMLRendered(context.Background(), &proto.ContentFile{Path: path, Content: content})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

// OnPostBuild is called after the build process is finished.
func (m *GRPCClient) OnPostBuild() error {
	_, err := m.Client.OnPostBuild(context.Background(), &proto.PostBuildRequest{})
	return err
}

// GRPCServer is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// Impl is the real implementation of the plugin.
	Impl Plugin
	proto.UnimplementedPluginServer
}

// OnPreBuild is called before the build process starts.
func (m *GRPCServer) OnPreBuild(ctx context.Context, req *proto.PreBuildRequest) (*proto.PreBuildResponse, error) {
	return &proto.PreBuildResponse{}, m.Impl.OnPreBuild()
}

// OnConfigLoaded is called after the configuration is loaded.
func (m *GRPCServer) OnConfigLoaded(ctx context.Context, req *proto.ConfigLoadedRequest) (*proto.ConfigLoadedResponse, error) {
	config, err := m.Impl.OnConfigLoaded([]byte(req.ConfigJson))
	if err != nil {
		return nil, err
	}
	return &proto.ConfigLoadedResponse{ConfigJson: string(config)}, nil
}

// OnPublicAssetsCopied is called after the public assets are copied.
func (m *GRPCServer) OnPublicAssetsCopied(ctx context.Context, req *proto.PublicAssetsCopiedRequest) (*proto.PublicAssetsCopiedResponse, error) {
	return &proto.PublicAssetsCopiedResponse{}, m.Impl.OnPublicAssetsCopied()
}

// OnContentLoaded is called after a content file is loaded.
func (m *GRPCServer) OnContentLoaded(ctx context.Context, req *proto.ContentFile) (*proto.ContentFile, error) {
	content, err := m.Impl.OnContentLoaded(req.Path, req.Content)
	if err != nil {
		return nil, err
	}
	return &proto.ContentFile{Path: req.Path, Content: content}, nil
}

// OnContentRender is called after a content file is rendered.
func (m *GRPCServer) OnContentRender(ctx context.Context, req *proto.ContentFile) (*proto.ContentFile, error) {
	content, err := m.Impl.OnContentRender(req.Path, req.Content)
	if err != nil {
		return nil, err
	}
	return &proto.ContentFile{Path: req.Path, Content: content}, nil
}

// OnHTMLRendered is called after the HTML is rendered.
func (m *GRPCServer) OnHTMLRendered(ctx context.Context, req *proto.ContentFile) (*proto.ContentFile, error) {
	content, err := m.Impl.OnHTMLRendered(req.Path, req.Content)
	if err != nil {
		return nil, err
	}
	return &proto.ContentFile{Path: req.Path, Content: content}, nil
}

// OnPostBuild is called after the build process is finished.
func (m *GRPCServer) OnPostBuild(ctx context.Context, req *proto.PostBuildRequest) (*proto.PostBuildResponse, error) {
	return &proto.PostBuildResponse{}, m.Impl.OnPostBuild()
}
