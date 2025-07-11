package plugins_test

import (
	"net"
	"testing"

	"github.com/Bitlatte/evoke/pkg/plugins"
	"github.com/Bitlatte/evoke/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// mockPlugin is a mock implementation of the Plugin interface.
type mockPlugin struct{}

func (m *mockPlugin) Name() string      { return "mock" }
func (m *mockPlugin) OnPreBuild() error { return nil }
func (m *mockPlugin) OnConfigLoaded(config []byte) ([]byte, error) {
	return config, nil
}
func (m *mockPlugin) OnPublicAssetsCopied() error { return nil }
func (m *mockPlugin) OnContentLoaded(path string, content []byte) ([]byte, error) {
	return content, nil
}
func (m *mockPlugin) OnContentRender(path string, content []byte) ([]byte, error) {
	return content, nil
}
func (m *mockPlugin) OnHTMLRendered(path string, content []byte) ([]byte, error) {
	return content, nil
}
func (m *mockPlugin) OnPostBuild() error { return nil }
func (m *mockPlugin) RegisterPipelines() ([]*proto.Pipeline, error) {
	return []*proto.Pipeline{
		{
			Name:       "test",
			Extensions: []string{".test"},
		},
	}, nil
}
func (m *mockPlugin) ProcessAsset(asset *proto.Asset) (*proto.Asset, error) {
	return asset, nil
}

func TestPlugin(t *testing.T) {
	// Create a mock server
	server := grpc.NewServer()
	proto.RegisterPluginServer(server, &plugins.GRPCServer{Impl: &mockPlugin{}})

	// Create a listener
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Serve the server in a goroutine
	go server.Serve(lis)

	// Create a client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer conn.Close()

	// Create the plugin client
	client := &plugins.EvokeGRPCClient{Client: proto.NewPluginClient(conn)}

	// Test the plugin methods
	if err := client.OnPreBuild(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func BenchmarkPlugin(b *testing.B) {
	// Create a mock server
	server := grpc.NewServer()
	proto.RegisterPluginServer(server, &plugins.GRPCServer{Impl: &mockPlugin{}})

	// Create a listener
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	// Serve the server in a goroutine
	go server.Serve(lis)

	// Create a client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Fatalf("err: %s", err)
	}
	defer conn.Close()

	// Create the plugin client
	client := &plugins.EvokeGRPCClient{Client: proto.NewPluginClient(conn)}
	content := make([]byte, 1024*10) // 10KB

	b.ResetTimer()
	b.ReportAllocs()

	// Run the benchmark on a method that transfers data
	for i := 0; i < b.N; i++ {
		_, err := client.OnContentLoaded("path/to/content.md", content)
		if err != nil {
			b.Fatalf("err: %s", err)
		}
	}
}
