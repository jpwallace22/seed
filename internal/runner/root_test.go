package runner

import (
	"errors"
	"os"
	"testing"

	"github.com/jpwallace22/seed/cmd/flags"
	"github.com/jpwallace22/seed/internal/ctx"
	mocklogger "github.com/jpwallace22/seed/pkg/logger/mock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/* ******************************************
*               Define mocks
* *******************************************/

type MockClipboard struct {
	mock.Mock
}

func (m *MockClipboard) CopyText(text string) error {
	args := m.Called(text)
	return args.Error(0)
}

func (m *MockClipboard) PasteText() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type MockParser struct {
	mock.Mock
}

func (m *MockParser) ParseTree(tree string) error {
	args := m.Called(tree)
	return args.Error(0)
}

func buildTestRunner(testFlags flags.RootFlags) (*RootRunner, *MockClipboard, *MockParser) {
	mockLogger := mocklogger.New()
	mockClipboard := new(MockClipboard)
	mockParser := new(MockParser)
	mockCmd := &cobra.Command{
		Use: "test",
	}

	testCtx := &ctx.SeedContext{
		Logger: mockLogger,
		Cobra:  mockCmd,
		Flags: flags.Flags{
			Root: testFlags,
		},
	}

	runner := &RootRunner{
		ctx:       testCtx,
		clipboard: mockClipboard,
		parser:    mockParser,
	}

	return runner, mockClipboard, mockParser
}

/* ******************************************
*               Tests
*********************************************/
func TestClipboardOperations(t *testing.T) {
	tests := []struct {
		clipError     error
		name          string
		clipContent   string
		errorContains string
		flags         flags.RootFlags
		args          []string
		expectError   bool
	}{
		{
			name: "successful clipboard parse",
			flags: flags.RootFlags{
				FromClipboard: true,
			},
			clipContent: "test-content",
			expectError: false,
		},
		{
			name: "clipboard read error",
			flags: flags.RootFlags{
				FromClipboard: true,
			},
			clipError:     errors.New("clipboard error"),
			expectError:   true,
			errorContains: "clipboard read error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, mockClipboard, mockParser := buildTestRunner(tt.flags)

			mockClipboard.On("PasteText").Return(tt.clipContent, tt.clipError)
			if !tt.expectError {
				mockParser.On("ParseTree", tt.clipContent).Return(nil)
			}

			err := runner.Run(tt.args)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockClipboard.AssertExpectations(t)
			mockParser.AssertExpectations(t)
		})
	}
}

func TestFileOperations(t *testing.T) {
	tests := []struct {
		name          string
		filePath      string
		fileContent   string
		errorContains string
		flags         flags.RootFlags
		args          []string
		expectError   bool
	}{
		{
			name: "successful file parse",
			flags: flags.RootFlags{
				FilePath: "test.txt",
			},
			fileContent: "test content",
			expectError: false,
		},
		{
			name: "file not found",
			flags: flags.RootFlags{
				FilePath: "nonexistent.txt",
			},
			expectError:   true,
			errorContains: "file read error",
		},
		{
			name: "empty file",
			flags: flags.RootFlags{
				FilePath: "empty.txt",
			},
			expectError:   true,
			errorContains: "file read error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, _, mockParser := buildTestRunner(tt.flags)

			if !tt.expectError && tt.fileContent != "" {
				err := os.WriteFile(tt.flags.FilePath, []byte(tt.fileContent), 0644)
				defer os.Remove(tt.flags.FilePath)
				assert.NoError(t, err)
				mockParser.On("ParseTree", tt.fileContent).Return(nil)
			}

			err := runner.Run(tt.args)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockParser.AssertExpectations(t)
		})
	}
}
