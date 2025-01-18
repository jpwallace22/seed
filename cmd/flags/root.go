package flags

import "fmt"

type RootFlags struct {
	FilePath      string
	Format        Format
	Silent        bool
	FromClipboard bool
}

type Format string

const (
	Tree Format = "tree"
	JSON Format = "json"
	YAML Format = "yaml"
)

func (f Format) String() string {
	return string(f)
}

func (f *Format) Set(value string) error {
	switch Format(value) {
	case Tree, JSON, YAML:
		*f = Format(value)
		return nil
	default:
		return fmt.Errorf("invalid format %q, must be one of: tree, json, yaml", value)
	}
}

func (f Format) Type() string {
	return "format"
}
