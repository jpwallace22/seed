package mock

import "github.com/stretchr/testify/mock"

type MockLogger struct {
	mock.Mock
	defaultBehavior bool
}

func New() *MockLogger {
	return NewWithBehavior(true) // By default, use no-op behavior
}

func NewWithBehavior(useDefault bool) *MockLogger {
	m := &MockLogger{
		defaultBehavior: useDefault,
	}
	if useDefault {
		m.On("Info", mock.Anything, mock.Anything).Return(nil)
		m.On("Warn", mock.Anything, mock.Anything).Return(nil)
		m.On("Error", mock.Anything, mock.Anything).Return(nil)
		m.On("Log", mock.Anything, mock.Anything).Return(nil)
		m.On("Success", mock.Anything, mock.Anything).Return(nil)
	}
	return m
}

// Helper method to clear default expectations
func (m *MockLogger) ClearDefaultExpectations() {
	m.ExpectedCalls = nil
	m.defaultBehavior = false
}

func (m *MockLogger) Info(msg string, v ...interface{})    { m.Called(msg, v) }
func (m *MockLogger) Warn(msg string, v ...interface{})    { m.Called(msg, v) }
func (m *MockLogger) Error(msg string, v ...interface{})   { m.Called(msg, v) }
func (m *MockLogger) Log(msg string, v ...interface{})     { m.Called(msg, v) }
func (m *MockLogger) Success(msg string, v ...interface{}) { m.Called(msg, v) }
