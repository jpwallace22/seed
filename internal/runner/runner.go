package runner

type Runner[T any] interface {
	Run(flags T, args []string) error
}

type Config struct {
	Silent bool
}
