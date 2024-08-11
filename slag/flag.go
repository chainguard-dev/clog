// Package slag provides a method for setting the log level from the command line.
//
//	func main() {
//		var level slag.Level
//		flag.Var(&level, "log-level", "log level")
//		flag.Parse()
//		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: &level})))
//	}
//
// See [./examples/logger](./examples/logger) for a full example.
//
// This allows the log level to be set from the command line:
//
//	$ ./myprogram -log-level=debug
//
// The slag.Level type is a wrapper around slog.Level that implements the flag.Value interface,
// as well as Cobra's pflag.Value interface.
//
//	func main() {
//		var level slag.Level
//		cmd := &cobra.Command{
//			Use: "myprogram",
//			...
//		}
//		cmd.PersistentFlags().Var(&level, "log-level", "log level")
//		cmd.Execute()
//	}
package slag

import "log/slog"

type Level slog.Level

func (l *Level) Set(s string) error {
	var ll slog.Level
	if err := ll.UnmarshalText([]byte(s)); err != nil {
		return err
	}
	*l = Level(ll)
	return nil
}
func (l *Level) String() string    { return slog.Level(*l).String() }
func (l *Level) Level() slog.Level { return slog.Level(*l) }

// Implements https://pkg.go.dev/github.com/spf13/pflag#Value
func (l *Level) Type() string { return "string" }
