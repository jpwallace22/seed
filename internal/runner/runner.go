package runner

type Runner interface {
	Run(args []string) error
}

type Config struct {
	Silent bool
}
