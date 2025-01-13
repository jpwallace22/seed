package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jpwallace22/seed/internal/ctx"
	"github.com/jpwallace22/seed/pkg/logger"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
	logger  *mockLogger
	parser  *Parser
	tempDir string
}

type mockLogger struct {
	InfoLogs  []string
	WarnLogs  []string
	ErrorLogs []string
}

func (m *mockLogger) Info(msg string, v ...interface{})         { m.InfoLogs = append(m.InfoLogs, msg) }
func (m *mockLogger) Warn(msg string, v ...interface{})         { m.WarnLogs = append(m.WarnLogs, msg) }
func (m *mockLogger) Error(msg string, v ...interface{})        { m.ErrorLogs = append(m.ErrorLogs, msg) }
func (m *mockLogger) Log(msg string, v ...interface{})          { m.ErrorLogs = append(m.ErrorLogs, msg) }
func (m *mockLogger) Success(msg string, v ...interface{})      { m.ErrorLogs = append(m.ErrorLogs, msg) }
func (m *mockLogger) WithField(key, value string) logger.Logger { return m }

func (s *ParserTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "parser_test_*")
	s.Require().NoError(err)
	s.Require().NoError(os.Chdir(s.tempDir))

	s.logger = &mockLogger{}
	testCtx := ctx.SeedContext{
		Logger: s.logger,
	}
	s.parser = New(testCtx)
}

func (s *ParserTestSuite) TearDownTestSuite() {
	os.RemoveAll(s.tempDir)
}

func (s *ParserTestSuite) TestEmptyInput() {
	s.Run("empty input should error", func() {
		err := s.parser.ParseTreeString("")
		s.Error(err, "Expected error for empty input")
	})
}

func (s *ParserTestSuite) TestTreePrefix() {
	input := `tree
root
├── file1.txt
└── file2.txt`

	expectedFiles := []string{
		"root/file1.txt",
		"root/file2.txt",
	}
	expectedDirs := []string{"root"}

	s.Run("create tree with prefix", func() {
		s.Require().NoError(s.parser.ParseTreeString(input))
		s.verifyStructure(expectedFiles, expectedDirs)
	})
}

func (s *ParserTestSuite) TestSimpleStructure() {
	input := `root
├── dir1
└── dir2
    └── file.txt`

	expectedFiles := []string{"root/dir2/file.txt"}
	expectedDirs := []string{"root", "root/dir1", "root/dir2"}

	s.Run("create directory structure", func() {
		s.Require().NoError(s.parser.ParseTreeString(input))
		s.verifyStructure(expectedFiles, expectedDirs)
	})
}

func (s *ParserTestSuite) TestStripsTrailingSlashes() {
	input := `root
├── dir1\
└── dir2/
    └── file.txt`

	expectedFiles := []string{"root/dir2/file.txt"}
	expectedDirs := []string{"root", "root/dir1", "root/dir2"}

	s.Run("create directory structure", func() {
		s.Require().NoError(s.parser.ParseTreeString(input))
		s.verifyStructure(expectedFiles, expectedDirs)
	})
}

func (s *ParserTestSuite) TestComplexStructure() {
	input := `project
├── src
│   ├── main.go
│   └── utils
│       └── helper.go
└── tests
    ├── main_test.go
    └── utils
        └── helper_test.go`

	expectedFiles := []string{
		"project/src/main.go",
		"project/src/utils/helper.go",
		"project/tests/main_test.go",
		"project/tests/utils/helper_test.go",
	}
	expectedDirs := []string{
		"project",
		"project/src",
		"project/src/utils",
		"project/tests",
		"project/tests/utils",
	}

	s.Run("create nested directory structure", func() {
		s.Require().NoError(s.parser.ParseTreeString(input))
		s.verifyStructure(expectedFiles, expectedDirs)
	})
}

func (s *ParserTestSuite) TestDotRoot() {
	input := `.
├── poopy
│   └── bar
│       ├── baz
│       │   └── boop.txt
│       └── foop.js
└── test
    ├── foo
    │   └── bar.jpg
    └── test.txt`

	expectedFiles := []string{
		"poopy/bar/baz/boop.txt",
		"poopy/bar/foop.js",
		"test/foo/bar.jpg",
		"test/test.txt",
	}
	expectedDirs := []string{
		"poopy",
		"poopy/bar",
		"poopy/bar/baz",
		"test",
		"test/foo",
	}

	s.Run("create structure with dot root", func() {
		s.Require().NoError(s.parser.ParseTreeString(input))
		s.verifyStructure(expectedFiles, expectedDirs)

		// Additional checks for correct nesting
		s.FileExists(filepath.Join(s.tempDir, "poopy/bar/baz/boop.txt"))
		s.FileExists(filepath.Join(s.tempDir, "test/foo/bar.jpg"))

		// Verify directory existence explicitly
		s.DirExists(filepath.Join(s.tempDir, "poopy/bar/baz"))
		s.DirExists(filepath.Join(s.tempDir, "test/foo"))
	})
}

func (s *ParserTestSuite) TestRealWorldExample() {
	input := `poop
├── poopy
└── test
    ├── foo
    │   └── bar.jpg
    └── test.txt`

	expectedFiles := []string{
		"poop/test/foo/bar.jpg",
		"poop/test/test.txt",
	}
	expectedDirs := []string{
		"poop",
		"poop/poopy",
		"poop/test",
		"poop/test/foo",
	}

	s.Run("create real world structure", func() {
		s.Require().NoError(s.parser.ParseTreeString(input))
		s.verifyStructure(expectedFiles, expectedDirs)
	})
}

func (s *ParserTestSuite) TestDeepNesting() {
	input := `root
├── level1
│   ├── level2
│   │   ├── level3
│   │   │   └── deep.txt
│   │   └── file2.txt
│   └── file1.txt
└── sibling
    └── nephew.txt`

	expectedFiles := []string{
		"root/level1/level2/level3/deep.txt",
		"root/level1/level2/file2.txt",
		"root/level1/file1.txt",
		"root/sibling/nephew.txt",
	}
	expectedDirs := []string{
		"root",
		"root/level1",
		"root/level1/level2",
		"root/level1/level2/level3",
		"root/sibling",
	}

	s.Run("create deeply nested structure", func() {
		s.Require().NoError(s.parser.ParseTreeString(input))
		s.verifyStructure(expectedFiles, expectedDirs)

		// Verify specific deep nesting
		deepFile := filepath.Join(s.tempDir, "root/level1/level2/level3/deep.txt")
		s.FileExists(deepFile)

		// Verify all intermediate directories exist
		s.DirExists(filepath.Join(s.tempDir, "root/level1/level2/level3"))
	})
}

func (s *ParserTestSuite) TestMultipleSiblings() {
	input := `project
├── src
│   ├── file1.txt
│   ├── file2.txt
│   └── file3.txt
└── test
    ├── test1.txt
    ├── test2.txt
    └── test3.txt`

	expectedFiles := []string{
		"project/src/file1.txt",
		"project/src/file2.txt",
		"project/src/file3.txt",
		"project/test/test1.txt",
		"project/test/test2.txt",
		"project/test/test3.txt",
	}
	expectedDirs := []string{
		"project",
		"project/src",
		"project/test",
	}

	s.Run("create structure with multiple siblings", func() {
		s.Require().NoError(s.parser.ParseTreeString(input))
		s.verifyStructure(expectedFiles, expectedDirs)

		// Verify sibling files are in correct directories
		for _, file := range []string{"file1.txt", "file2.txt", "file3.txt"} {
			s.FileExists(filepath.Join(s.tempDir, "project/src", file))
		}
		for _, file := range []string{"test1.txt", "test2.txt", "test3.txt"} {
			s.FileExists(filepath.Join(s.tempDir, "project/test", file))
		}
	})
}

func (s *ParserTestSuite) TestGetDepth() {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Depth with tree characters and single indentation",
			input:    "│   └── file.txt",
			expected: 2,
		},
		{
			name:     "Depth with spaces only and single indentation",
			input:    "    └── file.txt",
			expected: 2,
		},
		{
			name:     "Depth with root tree character",
			input:    "├── dir",
			expected: 1,
		},
		{
			name:     "Depth with deeply indented line",
			input:    "        └── file.txt",
			expected: 3,
		},
		{
			name:     "Depth with no indentation",
			input:    "file.txt",
			expected: 0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := s.parser.getDepth(tt.input)
			s.Equal(tt.expected, result, "getDepth(%q) = %d, want %d", tt.input, result, tt.expected)
		})
	}
}

func (s *ParserTestSuite) TestExtractName() {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Extract file name from tree line",
			input:    "├── file.txt",
			expected: "file.txt",
		},
		{
			name:     "Extract file name with nested structure",
			input:    "│   └── file.txt",
			expected: "file.txt",
		},
		{
			name:     "Extract directory name from tree line",
			input:    "└── dir",
			expected: "dir",
		},
		{
			name:     "Extract plain file name without tree characters",
			input:    "file.txt",
			expected: "file.txt",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := s.parser.extractName(tt.input)
			s.Equal(tt.expected, result, "extractName(%q) = %q, want %q", tt.input, result, tt.expected)
		})
	}
}

func (s *ParserTestSuite) verifyStructure(expectedFiles, expectedDirs []string) {
	var actualFiles, actualDirs []string

	err := filepath.Walk(s.tempDir, func(path string, info os.FileInfo, err error) error {
		s.Require().NoError(err)
		if path == s.tempDir {
			return nil
		}

		relPath, err := filepath.Rel(s.tempDir, path)
		s.Require().NoError(err)

		// Normalize path for windows support
		relPath = filepath.ToSlash(relPath)

		if info.IsDir() {
			actualDirs = append(actualDirs, relPath)
		} else {
			actualFiles = append(actualFiles, relPath)
		}
		return nil
	})
	s.Require().NoError(err)

	s.ElementsMatch(expectedFiles, actualFiles, "Files created don't match expected")
	s.ElementsMatch(expectedDirs, actualDirs, "Directories created don't match expected")
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
