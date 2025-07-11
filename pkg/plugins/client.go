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
	return &GRPCClient{Client: proto.NewPluginClient(c)}, nil
}

// Plugin is the interface that all evoke plugins must implement.
type Plugin interface {
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
		p, ok := raw.(Plugin)
		if !ok {
			return err
		}

		plugins = append(plugins, p)

		return nil
	})

	return plugins, err
}
