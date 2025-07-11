package plugins

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Bitlatte/evoke/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "EVOKE_PLUGIN",
	MagicCookieValue: "1.0",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"evoke": &EvokePlugin{},
}

// EvokePlugin is the implementation of plugin.Plugin so we can serve/consume plugins.
type EvokePlugin struct {
	plugin.Plugin
	// Impl is the concrete implementation of the plugin.
	Impl Plugin
}

// GRPCServer registers the plugin with the gRPC server.
func (p *EvokePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

// GRPCClient returns the gRPC client for the plugin.
func (p *EvokePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &EvokeGRPCClient{Client: proto.NewPluginClient(c)}, nil
}

// Plugin is the interface that all evoke plugins must implement.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string
	// OnPreBuild is called before the build process starts.
	OnPreBuild() error
	// OnConfigLoaded is called after the configuration is loaded.
	OnConfigLoaded(config []byte) ([]byte, error)
	// OnPublicAssetsCopied is called after the public assets are copied.
	OnPublicAssetsCopied() error
	// OnContentLoaded is called after a content file is loaded.
	OnContentLoaded(path string, content []byte) ([]byte, error)
	// OnContentRender is called after a content file is rendered.
	OnContentRender(path string, content []byte) ([]byte, error)
	// OnHTMLRendered is called after the HTML is rendered.
	OnHTMLRendered(path string, content []byte) ([]byte, error)
	// OnPostBuild is called after the build process is finished.
	OnPostBuild() error
	// RegisterPipelines is called to register custom pipelines.
	RegisterPipelines() ([]*proto.Pipeline, error)
	// ProcessAsset is called to process an asset with a custom pipeline.
	ProcessAsset(asset *proto.Asset) (*proto.Asset, error)
}

// EvokeGRPCClient is an implementation of Plugin that talks over RPC.
type EvokeGRPCClient struct {
	Client proto.PluginClient
	name   string
}

// Name returns the name of the plugin.
func (m *EvokeGRPCClient) Name() string {
	return m.name
}

// OnPreBuild is called before the build process starts.
func (m *EvokeGRPCClient) OnPreBuild() error {
	_, err := m.Client.OnPreBuild(context.Background(), &proto.PreBuildRequest{})
	return err
}

// OnConfigLoaded is called after the configuration is loaded.
func (m *EvokeGRPCClient) OnConfigLoaded(config []byte) ([]byte, error) {
	resp, err := m.Client.OnConfigLoaded(context.Background(), &proto.ConfigLoadedRequest{
		ConfigJson: string(config),
	})
	if err != nil {
		return nil, err
	}
	return []byte(resp.ConfigJson), nil
}

// OnPublicAssetsCopied is called after the public assets are copied.
func (m *EvokeGRPCClient) OnPublicAssetsCopied() error {
	_, err := m.Client.OnPublicAssetsCopied(context.Background(), &proto.PublicAssetsCopiedRequest{})
	return err
}

// OnContentLoaded is called after a content file is loaded.
func (m *EvokeGRPCClient) OnContentLoaded(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnContentLoaded(context.Background(), &proto.ContentFile{
		Path:    path,
		Content: content,
	})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

// OnContentRender is called after a content file is rendered.
func (m *EvokeGRPCClient) OnContentRender(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnContentRender(context.Background(), &proto.ContentFile{
		Path:    path,
		Content: content,
	})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

// OnHTMLRendered is called after the HTML is rendered.
func (m *EvokeGRPCClient) OnHTMLRendered(path string, content []byte) ([]byte, error) {
	resp, err := m.Client.OnHTMLRendered(context.Background(), &proto.ContentFile{
		Path:    path,
		Content: content,
	})
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

// OnPostBuild is called after the build process is finished.
func (m *EvokeGRPCClient) OnPostBuild() error {
	_, err := m.Client.OnPostBuild(context.Background(), &proto.PostBuildRequest{})
	return err
}

// RegisterPipelines is called to register custom pipelines.
func (m *EvokeGRPCClient) RegisterPipelines() ([]*proto.Pipeline, error) {
	resp, err := m.Client.RegisterPipelines(context.Background(), &proto.RegisterPipelinesRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Pipelines, nil
}

// ProcessAsset is called to process an asset with a custom pipeline.
func (m *EvokeGRPCClient) ProcessAsset(asset *proto.Asset) (*proto.Asset, error) {
	return m.Client.ProcessAsset(context.Background(), asset)
}

// LoadPlugins loads all the plugins in the plugins directory.
func LoadPlugins() ([]Plugin, error) {
	var plugins []Plugin
	// We're going to walk the plugins directory and look for executable files
	err := filepath.Walk("plugins", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If it's a directory, skip it
		if info.IsDir() {
			return nil
		}

		// If it's not executable, skip it
		if info.Mode()&0111 == 0 {
			return nil
		}

		// Create a new plugin client
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: Handshake,
			Plugins:         PluginMap,
			Cmd:             exec.Command(path),
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:  "plugin",
				Level: hclog.Error,
			}),
		})

		// Connect to the plugin
		rpcClient, err := client.Client()
		if err != nil {
			return err
		}

		// Request the plugin
		raw, err := rpcClient.Dispense("evoke")
		if err != nil {
			return err
		}

		// Assert that the plugin is the correct type
		p, ok := raw.(*EvokeGRPCClient)
		if !ok {
			return err
		}

		p.name = filepath.Base(path)
		plugins = append(plugins, p)

		return nil
	})

	return plugins, err
}
