package runner

import (
	"errors"
	"testing"

	"github.com/jpwallace22/seed/internal/ctx"
	mocklogger "github.com/jpwallace22/seed/pkg/logger/mock"
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

	testCtx := &ctx.SeedContext{
		Logger: mockLogger,
		GlobalFlags: ctx.GlobalFlags{
			Silent: silent,
		},
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
		name  string
		flags RootFlags
	}{
		{
			name: "creates runner with all flags disabled",
			flags: RootFlags{
				Silent:        false,
				FromClipboard: false,
			},
		},
		{
			name: "creates runner with silent mode enabled",
			flags: RootFlags{
				Silent:        true,
				FromClipboard: false,
			},
		},
		{
			name: "creates runner with clipboard mode enabled",
			flags: RootFlags{
				Silent:        false,
				FromClipboard: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewRootRunner(tt.flags)

			assert.NotNil(t, runner)
			assert.NotNil(t, runner.ctx)
			assert.NotNil(t, runner.clipboard)
			assert.Equal(t, tt.flags.Silent, runner.ctx.GlobalFlags.Silent)
		})
	}
}

func TestRunnerRun(t *testing.T) {
	tests := []struct {
		clipError     error
		name          string
		clipContent   string
		errorContains string
		args          []string
		flags         RootFlags
		expectError   bool
	}{
		{
			name: "successful clipboard read",
			flags: RootFlags{
				FromClipboard: true,
				Silent:        false,
			},
			clipContent: "test-content",
			clipError:   nil,
			expectError: false,
		},
		{
			name: "clipboard read error",
			flags: RootFlags{
				FromClipboard: true,
				Silent:        false,
			},
			clipContent:   "",
			clipError:     errors.New("clipboard error"),
			expectError:   true,
			errorContains: "unable to access clipboard contents",
		},
		{
			name: "no input source provided",
			flags: RootFlags{
				FromClipboard: false,
				Silent:        false,
			},
			args:          []string{},
			expectError:   true,
			errorContains: "no seeds provided",
		},
		{
			name: "file path provided",
			flags: RootFlags{
				FromClipboard: false,
				Silent:        false,
			},
			args:        []string{"test.txt"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, mockClipboard, mockParser := buildTestRunner(tt.flags.Silent)

			if tt.flags.FromClipboard {
				mockClipboard.On("PasteText").Return(tt.clipContent, tt.clipError)
				if tt.clipError == nil {
					mockParser.On("ParseTreeString", tt.clipContent).Return(nil)
				}
			}

			if len(tt.args) > 0 {
				mockParser.On("ParseTreeString", tt.args[0]).Return(nil)
			}
			err := runner.Run(tt.flags.FromClipboard, tt.args)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockParser.AssertExpectations(t)
			mockClipboard.AssertExpectations(t)
		})
	}
}

func TestGetClipboardContent(t *testing.T) {
	tests := []struct {
		err         error
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "successful clipboard read",
			content:     "test content",
			err:         nil,
			expectError: false,
		},
		{
			name:        "clipboard read error",
			content:     "",
			err:         errors.New("mock clipboard error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, mockClipboard, _ := buildTestRunner(true)

			mockClipboard.On("PasteText").Return(tt.content, tt.err)
			content, err := runner.getClipboardContent()

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, content)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.content, content)
			}

			// Verify all mock expectations were met
			mockClipboard.AssertExpectations(t)
		})
	}
}
