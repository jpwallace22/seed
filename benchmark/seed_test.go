package benchmark

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// getBinaryPath returns the path to the seed binary
func getBinaryPath() (string, error) {
	// Get the directory where the test file is located
	_, filename, _, _ := runtime.Caller(0)
	benchmarkDir := filepath.Dir(filename)

	// Navigate up to project root and then to bin/seed
	binPath := filepath.Join(benchmarkDir, "..", "bin", "seed")

	// Get absolute path and verify binary exists
	absPath, err := filepath.Abs(binPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("binary not found at %s: %v", absPath, err)
	}

	return absPath, nil
}

func generateLargeTree() string {
	// Generate a large tree structure programmatically
	tree := "large-project\n"
	for i := 0; i < 10; i++ {
		tree += fmt.Sprintf("├── module%d\n", i)
		for j := 0; j < 5; j++ {
			tree += fmt.Sprintf("│   ├── submodule%d\n", j)
			for k := 0; k < 4; k++ {
				tree += fmt.Sprintf("│   │   ├── file%d.ts\n", k)
			}
		}
	}
	return tree
}

func TestMain(m *testing.M) {
	// Ensure we're in the benchmark directory
	if err := os.Chdir(filepath.Dir(os.Args[0])); err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// Benchmark ASCII tree parsing with different sizes
func BenchmarkParseASCIITree(b *testing.B) {
	binaryPath, err := getBinaryPath()
	if err != nil {
		b.Fatal(err)
	}

	tests := []struct {
		name string
		tree string
	}{
		{"100 Nodes", generateTreeStructure(SmallSize)},
		{"500 Nodes", generateTreeStructure(MediumSize)},
		{"1000 Nodes", generateTreeStructure(LargeSize)},
		{"5000 Nodes", generateTreeStructure(ExtraLargeSize)},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			tmpDir, err := os.MkdirTemp("", "seed-benchmark-*")
			if err != nil {
				b.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			tmpFile := filepath.Join(tmpDir, "tree.txt")
			if err := os.WriteFile(tmpFile, []byte(tt.tree), 0644); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cmd := exec.Command(binaryPath, "-f", tmpFile)
				cmd.Dir = tmpDir
				output, err := cmd.CombinedOutput()
				if err != nil {
					b.Fatalf("command failed: %v\nOutput: %s", err, output)
				}
			}
		})
	}
}

// Benchmark JSON parsing with different sizes
func BenchmarkParseJSON(b *testing.B) {
	binaryPath, err := getBinaryPath()
	if err != nil {
		b.Fatal(err)
	}

	tests := []struct {
		name string
		json []interface{}
	}{
		{"100 Nodes", generateJSONStructure(SmallSize)},
		{"500 Nodes", generateJSONStructure(MediumSize)},
		{"1000 Nodes", generateJSONStructure(LargeSize)},
		{"5000 Nodes", generateJSONStructure(ExtraLargeSize)},
	}

	for _, tt := range tests {
		name := tt.name
		b.Run(name, func(b *testing.B) {
			tmpDir, err := os.MkdirTemp("", "seed-benchmark-*")
			if err != nil {
				b.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			jsonData, err := json.Marshal(tt.json)
			if err != nil {
				b.Fatal(err)
			}

			tmpFile := filepath.Join(tmpDir, "structure.json")
			if err := os.WriteFile(tmpFile, jsonData, 0644); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cmd := exec.Command(binaryPath, "-F", "json", "-f", tmpFile)
				cmd.Dir = tmpDir
				output, err := cmd.CombinedOutput()
				if err != nil {
					b.Fatalf("command failed: %v\nOutput: %s", err, string(output))
				}
			}
		})
	}
}

// Benchmark string input vs file input
func BenchmarkInputMethods(b *testing.B) {
	binaryPath, err := getBinaryPath()
	if err != nil {
		b.Fatal(err)
	}

	tree := generateTreeStructure(MediumSize)

	b.Run("StringInput - 500 Nodes", func(b *testing.B) {
		tmpDir, err := os.MkdirTemp("", "seed-benchmark-*")
		if err != nil {
			b.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cmd := exec.Command(binaryPath, tree)
			cmd.Dir = tmpDir
			if err := cmd.Run(); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("FileInput - 500 nodes", func(b *testing.B) {
		tmpDir, err := os.MkdirTemp("", "seed-benchmark-*")
		if err != nil {
			b.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		tmpFile := filepath.Join(tmpDir, "tree.txt")
		if err := os.WriteFile(tmpFile, []byte(tree), 0644); err != nil {
			b.Fatal(err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cmd := exec.Command(binaryPath, "-f", tmpFile)
			cmd.Dir = tmpDir
			if err := cmd.Run(); err != nil {
				b.Fatal(err)
			}
		}
	})
}
