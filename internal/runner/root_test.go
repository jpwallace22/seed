package runner

import (
	"errors"
	"os"
	"testing"

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

func (m *MockParser) ParseTreeString(tree string) error {
	args := m.Called(tree)
	return args.Error(0)
}

func buildTestRunner(silent bool) (*RootRunner, *MockClipboard, *MockParser) {
	mockLogger := mocklogger.New()
	mockClipboard := new(MockClipboard)
	mockParser := new(MockParser)
	mockCmd := &cobra.Command{
		Use: "test",
	}

	testCtx := &ctx.SeedContext{
		Logger: mockLogger,
		GlobalFlags: ctx.GlobalFlags{
			Silent: silent,
		},
		Cobra: mockCmd,
	}

	runner := &RootRunner{
		ctx:       *testCtx,
		clipboard: mockClipboard,
		parser:    mockParser,
	}

	return runner, mockClipboard, mockParser
}

/* ******************************************
*               Tests
*********************************************/
func TestNewRunner(t *testing.T) {
	tests := []struct {
		name   string
		flags  RootFlags
		config Config
	}{
		{
			name:   "creates runner with all flags disabled",
			flags:  RootFlags{},
			config: Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use: "test",
			}
			runner := NewRootRunner(cmd, tt.config.Silent)
			rootRunner, ok := runner.(*RootRunner)

			assert.True(t, ok)
			assert.NotNil(t, runner)
			assert.NotNil(t, rootRunner.ctx)
			assert.NotNil(t, rootRunner.clipboard)
			assert.Equal(t, tt.config.Silent, rootRunner.ctx.GlobalFlags.Silent)
		})
	}
}

func TestClipboardOperations(t *testing.T) {
	tests := []struct {
		clipError     error
		name          string
		clipContent   string
		errorContains string
		flags         RootFlags
		args          []string
		expectError   bool
	}{
		{
			name: "successful clipboard parse",
			flags: RootFlags{
				FromClipboard: true,
			},
			clipContent: "test-content",
			expectError: false,
		},
		{
			name: "clipboard read error",
			flags: RootFlags{
				FromClipboard: true,
			},
			clipError:     errors.New("clipboard error"),
			expectError:   true,
			errorContains: "clipboard read error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, mockClipboard, mockParser := buildTestRunner(true)

			mockClipboard.On("PasteText").Return(tt.clipContent, tt.clipError)
			if !tt.expectError {
				mockParser.On("ParseTreeString", tt.clipContent).Return(nil)
			}

			err := runner.Run(tt.flags, tt.args)

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
		flags         RootFlags
		args          []string
		expectError   bool
	}{
		{
			name: "successful file parse",
			flags: RootFlags{
				FilePath: "test.txt",
			},
			fileContent: "test content",
			expectError: false,
		},
		{
			name: "file not found",
			flags: RootFlags{
				FilePath: "nonexistent.txt",
			},
			expectError:   true,
			errorContains: "file read error",
		},
		{
			name: "empty file",
			flags: RootFlags{
				FilePath: "empty.txt",
			},
			expectError:   true,
			errorContains: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, _, mockParser := buildTestRunner(true)

			if !tt.expectError && tt.fileContent != "" {
				err := os.WriteFile(tt.flags.FilePath, []byte(tt.fileContent), 0644)
				defer os.Remove(tt.flags.FilePath)
				assert.NoError(t, err)
				mockParser.On("ParseTreeString", tt.fileContent).Return(nil)
			}

			err := runner.Run(tt.flags, tt.args)

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
