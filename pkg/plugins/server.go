package plugins

import (
	"context"

	"github.com/Bitlatte/evoke/proto"
)

// GRPCClient is an implementation of Plugin that talks over RPC.
type GRPCClient struct {
	Client proto.PluginClient
}

func (m *GRPCClient) OnPreBuild() error {
	_, err := m.Client.OnPreBuild(context.Background(), &proto.PreBuildRequest{})
	return err
}

func (m *GRPCClient) OnConfigLoaded(config []byte) ([]byte, error) {
	resp, err := m.Client.OnConfigLoaded(context.Background(), &proto.ConfigLoadedRequest{ConfigJson: string(config)})
	if err != nil {
		return nil, err
	}
	return []byte(resp.ConfigJson), nil
}

func (m *GRPCClient) OnPublicAssetsCopied() error {
	_, err := m.Client.OnPublicAssetsCopied(context.Background(), &proto.PublicAssetsCopiedRequest{})
	return err
}

func (m *GRPCClient) OnContentLoaded(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnContentLoaded(context.Background(), &proto.ContentFile{Path: path, Content: content})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

func (m *GRPCClient) OnContentRender(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnContentRender(context.Background(), &proto.ContentFile{Path: path, Content: content})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

func (m *GRPCClient) OnHTMLRendered(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnHTMLRendered(context.Background(), &proto.ContentFile{Path: path, Content: content})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

func (m *GRPCClient) OnPostBuild() error {
	_, err := m.Client.OnPostBuild(context.Background(), &proto.PostBuildRequest{})
	return err
}

// GRPCServer is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Plugin
	proto.UnimplementedPluginServer
}

func (m *GRPCServer) OnPreBuild(ctx context.Context, req *proto.PreBuildRequest) (*proto.PreBuildResponse, error) {
	return &proto.PreBuildResponse{}, m.Impl.OnPreBuild()
}

func (m *GRPCServer) OnConfigLoaded(ctx context.Context, req *proto.ConfigLoadedRequest) (*proto.ConfigLoadedResponse, error) {
	config, err := m.Impl.OnConfigLoaded([]byte(req.ConfigJson))
	if err != nil {
		return nil, err
	}
	return &proto.ConfigLoadedResponse{ConfigJson: string(config)}, nil
}

func (m *GRPCServer) OnPublicAssetsCopied(ctx context.Context, req *proto.PublicAssetsCopiedRequest) (*proto.PublicAssetsCopiedResponse, error) {
	return &proto.PublicAssetsCopiedResponse{}, m.Impl.OnPublicAssetsCopied()
}

func (m *GRPCServer) OnContentLoaded(ctx context.Context, req *proto.ContentFile) (*proto.ContentFile, error) {
	content, err := m.Impl.OnContentLoaded(req.Path, req.Content)
	if err != nil {
		return nil, err
	}
	return &proto.ContentFile{Path: req.Path, Content: content}, nil
}

func (m *GRPCServer) OnContentRender(ctx context.Context, req *proto.ContentFile) (*proto.ContentFile, error) {
	content, err := m.Impl.OnContentRender(req.Path, req.Content)
	if err != nil {
		return nil, err
	}
	return &proto.ContentFile{Path: req.Path, Content: content}, nil
}

func (m *GRPCServer) OnHTMLRendered(ctx context.Context, req *proto.ContentFile) (*proto.ContentFile, error) {
	content, err := m.Impl.OnHTMLRendered(req.Path, req.Content)
	if err != nil {
		return nil, err
	}
	return &proto.ContentFile{Path: req.Path, Content: content}, nil
}

func (m *GRPCServer) OnPostBuild(ctx context.Context, req *proto.PostBuildRequest) (*proto.PostBuildResponse, error) {
	return &proto.PostBuildResponse{}, m.Impl.OnPostBuild()
}
