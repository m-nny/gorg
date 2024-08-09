package logutils

import (
	"flag"
	"log/slog"
)

var _ flag.Value = (*levelValue)(nil)

type levelValue slog.Level

func newLevelValue(val slog.Level, p *slog.Level) *levelValue {
	*p = val
	return (*levelValue)(p)
}

func (l *levelValue) Set(s string) error {
	v := slog.Level(0)
	err := v.UnmarshalText([]byte(s))
	*l = levelValue(v)
	return err
}

func (l *levelValue) String() string {
	return slog.Level(*l).String()
}

func LevelVar(p *slog.Level, name string, value slog.Level, usage string) {
	flag.Var(newLevelValue(value, p), name, usage)
}

func Level(name string, value slog.Level, usage string) *slog.Level {
	p := new(slog.Level)
	LevelVar(p, name, value, usage)
	return p
}
