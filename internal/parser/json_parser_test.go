package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jpwallace22/seed/internal/ctx"
	logMock "github.com/jpwallace22/seed/pkg/logger/mock"
	"github.com/stretchr/testify/suite"
)

type JsonTestSuite struct {
	suite.Suite
	logger  *logMock.MockLogger
	parser  Parser
	tempDir string
}

func (s *JsonTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "parser_test_*")
	s.Require().NoError(err)
	s.Require().NoError(os.Chdir(s.tempDir))

	s.logger = logMock.New()
	testCtx := ctx.SeedContext{
		Logger: s.logger,
	}
	s.parser = NewJSONParser(testCtx)
}

func (s *JsonTestSuite) TearDownTestSuite() {
	os.RemoveAll(s.tempDir)
}

func (s *JsonTestSuite) TestEmptyInput() {
	s.Run("empty input should error", func() {
		err := s.parser.ParseTree("")
		s.Error(err, "Expected error for empty input")
	})
}

func (s *JsonTestSuite) TestInvalidJSON() {
	s.Run("invalid JSON should error", func() {
		err := s.parser.ParseTree("not json")
		s.Error(err, "Expected error for invalid JSON")
	})

	s.Run("incomplete JSON should error", func() {
		err := s.parser.ParseTree(`[{"type":"directory","name":"root"`)
		s.Error(err, "Expected error for incomplete JSON")
	})
}

func (s *JsonTestSuite) TestSimpleStructure() {
	input := `[
		{"type":"directory","name":"root","contents":[
			{"type":"directory","name":"dir1"},
			{"type":"directory","name":"dir2","contents":[
				{"type":"file","name":"file.txt"}
			]}
		]},
		{"type":"report","directories":3,"files":1}
	]`

	expectedFiles := []string{"root/dir2/file.txt"}
	expectedDirs := []string{"root", "root/dir1", "root/dir2"}

	s.Run("create directory structure", func() {
		s.Require().NoError(s.parser.ParseTree(input))
		s.verifyStructure(expectedFiles, expectedDirs)
	})
}

func (s *JsonTestSuite) TestComplexStructure() {
	input := `[
		{"type":"directory","name":"project","contents":[
			{"type":"directory","name":"src","contents":[
				{"type":"file","name":"main.go"},
				{"type":"directory","name":"utils","contents":[
					{"type":"file","name":"helper.go"}
				]}
			]},
			{"type":"directory","name":"tests","contents":[
				{"type":"file","name":"main_test.go"},
				{"type":"directory","name":"utils","contents":[
					{"type":"file","name":"helper_test.go"}
				]}
			]}
		]},
		{"type":"report","directories":5,"files":4}
	]`

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
		s.Require().NoError(s.parser.ParseTree(input))
		s.verifyStructure(expectedFiles, expectedDirs)
	})
}

func (s *JsonTestSuite) TestDeepNesting() {
	input := `[
		{"type":"directory","name":"root","contents":[
			{"type":"directory","name":"level1","contents":[
				{"type":"directory","name":"level2","contents":[
					{"type":"directory","name":"level3","contents":[
						{"type":"file","name":"deep.txt"}
					]},
					{"type":"file","name":"file2.txt"}
				]},
				{"type":"file","name":"file1.txt"}
			]},
			{"type":"directory","name":"sibling","contents":[
				{"type":"file","name":"nephew.txt"}
			]}
		]},
		{"type":"report","directories":5,"files":4}
	]`

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
		s.Require().NoError(s.parser.ParseTree(input))
		s.verifyStructure(expectedFiles, expectedDirs)

		deepFile := filepath.Join(s.tempDir, "root/level1/level2/level3/deep.txt")
		s.FileExists(deepFile)
		s.DirExists(filepath.Join(s.tempDir, "root/level1/level2/level3"))
	})
}

func (s *JsonTestSuite) TestMultipleSiblings() {
	input := `[
		{"type":"directory","name":"project","contents":[
			{"type":"directory","name":"src","contents":[
				{"type":"file","name":"file1.txt"},
				{"type":"file","name":"file2.txt"},
				{"type":"file","name":"file3.txt"}
			]},
			{"type":"directory","name":"test","contents":[
				{"type":"file","name":"test1.txt"},
				{"type":"file","name":"test2.txt"},
				{"type":"file","name":"test3.txt"}
			]}
		]},
		{"type":"report","directories":3,"files":6}
	]`

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
		s.Require().NoError(s.parser.ParseTree(input))
		s.verifyStructure(expectedFiles, expectedDirs)

		for _, file := range []string{"file1.txt", "file2.txt", "file3.txt"} {
			s.FileExists(filepath.Join(s.tempDir, "project/src", file))
		}
		for _, file := range []string{"test1.txt", "test2.txt", "test3.txt"} {
			s.FileExists(filepath.Join(s.tempDir, "project/test", file))
		}
	})
}

func (s *JsonTestSuite) TestReportValidation() {
	s.Run("incorrect file count should error", func() {
		input := `[
			{"type":"directory","name":"root","contents":[
				{"type":"file","name":"file1.txt"}
			]},
			{"type":"report","directories":1,"files":2}
		]`
		err := s.parser.ParseTree(input)
		s.Error(err, "Expected error for incorrect file count in report")
	})

	s.Run("incorrect directory count should error", func() {
		input := `[
			{"type":"directory","name":"root","contents":[
				{"type":"directory","name":"dir1"}
			]},
			{"type":"report","directories":3,"files":0}
		]`
		err := s.parser.ParseTree(input)
		s.Error(err, "Expected error for incorrect directory count in report")
	})
}

func (s *JsonTestSuite) TestMissingFields() {
	s.Run("missing type should error", func() {
		input := `[
			{"name":"root","contents":[
				{"type":"file","name":"file1.txt"}
			]}
		]`
		err := s.parser.ParseTree(input)
		s.Error(err, "Expected error for missing type field")
	})

	s.Run("missing name should error", func() {
		input := `[
			{"type":"directory","contents":[
				{"type":"file","name":"file1.txt"}
			]}
		]`
		err := s.parser.ParseTree(input)
		s.Error(err, "Expected error for missing name field")
	})
}

func (s *JsonTestSuite) TestWithoutReport() {
	input := `[
		{"type":"directory","name":"root","contents":[
			{"type":"directory","name":"docs","contents":[
				{"type":"file","name":"readme.md"}
			]},
			{"type":"file","name":"config.json"}
		]}
	]`

	expectedFiles := []string{
		"root/docs/readme.md",
		"root/config.json",
	}
	expectedDirs := []string{
		"root",
		"root/docs",
	}

	s.Run("create structure without report section", func() {
		s.Require().NoError(s.parser.ParseTree(input))
		s.verifyStructure(expectedFiles, expectedDirs)

		// Verify specific files exist
		s.FileExists(filepath.Join(s.tempDir, "root/docs/readme.md"))
		s.FileExists(filepath.Join(s.tempDir, "root/config.json"))
	})
}

func (s *JsonTestSuite) verifyStructure(expectedFiles, expectedDirs []string) {
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

func TestJSONSuite(t *testing.T) {
	suite.Run(t, new(JsonTestSuite))
}
