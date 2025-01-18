package runner

type Runner[T any] interface {
	Run(args []string) error
}

type Config struct {
	Silent bool
}
