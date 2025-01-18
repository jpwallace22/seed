package flags

import "fmt"

type RootFlags struct {
	FilePath      string
	Format        Format
	Silent        bool
	FromClipboard bool
}

type Format string

var Formats = struct {
	Tree Format
	JSON Format
	YAML Format
}{
	Tree: "tree",
	JSON: "json",
	YAML: "yaml",
}

func (f Format) String() string {
	return string(f)
}

func (f *Format) Set(value string) error {
	switch Format(value) {
	case Formats.Tree, Formats.JSON, Formats.YAML:
		*f = Format(value)
		return nil
	default:
		return fmt.Errorf("invalid format %q, must be one of: tree, json, yaml", value)
	}
}

func (f Format) Type() string {
	return "format"
}
