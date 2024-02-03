package slag

import (
	"flag"
	"log/slog"
	"strings"
)

type Level slog.Level

func (l *Level) String() string { return slog.Level(*l).String() }
func (l *Level) Set(s string) error {
	var ll slog.Level
	if err := ll.UnmarshalJSON([]byte(strings.ToUpper(s))); err != nil {
		return err
	}
	*l = Level(ll)
	return nil
}
func (l *Level) Level() slog.Level { return slog.Level(*l) }

func LevelFlag(name string, value slog.Level, usage string) *Level {
	var l Level
	flag.CommandLine.Var(&l, name, usage)
	return &l
}
