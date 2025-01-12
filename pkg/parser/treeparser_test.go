package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jpwallace22/seed/pkg/ctx"
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

func (s *ParserTestSuite) TearDownTest() {
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

func (s *ParserTestSuite) verifyStructure(expectedFiles, expectedDirs []string) {
	var actualFiles, actualDirs []string

	err := filepath.Walk(s.tempDir, func(path string, info os.FileInfo, err error) error {
		s.Require().NoError(err)
		if path == s.tempDir {
			return nil
		}

		relPath, err := filepath.Rel(s.tempDir, path)
		s.Require().NoError(err)

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
