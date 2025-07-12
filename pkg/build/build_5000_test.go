package build_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/Bitlatte/evoke/pkg/build"
)

func BenchmarkBuild5000(b *testing.B) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "evoke-benchmark-5000")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to the temporary directory
	originalWd, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}
	err = os.Chdir(tmpDir)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Chdir(originalWd)

	// Generate a site with 5000 pages
	generateBenchmarkSite(b, 5000)

	b.ResetTimer()
	b.ReportAllocs()

	// Run the build
	for i := 0; i < b.N; i++ {
		// We must remove the dist directory on each iteration to get an accurate
		// measurement of a clean build.
		os.RemoveAll("dist")

		err = build.Build("dist", true, runtime.NumCPU())
		if err != nil {
			b.Fatal(err)
		}
	}
}
