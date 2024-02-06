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
