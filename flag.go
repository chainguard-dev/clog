package clog

import (
	"flag"
	"log/slog"
	"strings"
)

type Level struct{ l slog.Level }

func (l *Level) String() string     { return l.l.String() }
func (l *Level) Set(s string) error { return l.l.UnmarshalText([]byte(strings.ToUpper(s))) }
func (l *Level) Get() slog.Level    { return l.l }

func LevelFlag(name string, value slog.Level, usage string) *Level {
	l := &Level{value}
	flag.CommandLine.Var(l, name, usage)
	return l
}
