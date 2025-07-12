package dag

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Node represents a node in the dependency graph
type Node struct {
	Path         string
	Dependencies []*Node
}

// Graph represents the dependency graph
type Graph struct {
	Nodes map[string]*Node
}

// NewGraph creates a new dependency graph
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(path string) *Node {
	node := &Node{Path: path}
	g.Nodes[path] = node
	return node
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(from, to *Node) {
	from.Dependencies = append(from.Dependencies, to)
}

// BuildGraph builds the dependency graph for the given content directory
func BuildGraph(contentDir, partialsDir string) (*Graph, error) {
	graph := NewGraph()

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			graph.AddNode(path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, node := range graph.Nodes {
		dependencies, err := getDependencies(node.Path, partialsDir)
		if err != nil {
			return nil, err
		}

		for _, dependency := range dependencies {
			if _, ok := graph.Nodes[dependency]; !ok {
				graph.AddNode(dependency)
			}
			graph.AddEdge(node, graph.Nodes[dependency])
		}
	}

	return graph, nil
}

// getDependencies returns the dependencies for the given file
func getDependencies(path, partialsDir string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`{{\s*partial\s+"([^"]+)"\s*}}`)
	matches := r.FindAllStringSubmatch(string(content), -1)

	var dependencies []string
	for _, match := range matches {
		dependencies = append(dependencies, filepath.Join(partialsDir, match[1]))
	}

	return dependencies, nil
}

// GetDependents returns the dependents of the given node
func (g *Graph) GetDependents(node *Node) []*Node {
	var dependents []*Node
	for _, n := range g.Nodes {
		for _, d := range n.Dependencies {
			if d == node {
				dependents = append(dependents, n)
			}
		}
	}
	return dependents
}

// String returns a string representation of the graph
func (g *Graph) String() string {
	var b strings.Builder
	for _, node := range g.Nodes {
		b.WriteString(fmt.Sprintf("%s -> [", node.Path))
		for i, dep := range node.Dependencies {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(dep.Path)
		}
		b.WriteString("]\n")
	}
	return b.String()
}
