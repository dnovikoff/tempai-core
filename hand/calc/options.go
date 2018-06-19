package calc

import (
	"github.com/dnovikoff/tempai-core/compact"
)

type Form int

const (
	Regular Form = 1 << iota
	Pairs
	Kokushi
)

func (f Form) Check(x Form) bool {
	return (f & x) == x
}

type Options struct {
	Opened   int
	Used     compact.Instances
	Results  Results
	Declared Melds
	Forms    Form

	StartMelds Melds
}

type Option func(*Options)

func GetOptions(opts ...Option) *Options {
	r := Options{
		Forms: Kokushi | Pairs | Regular,
	}
	for _, opt := range opts {
		opt(&r)
	}
	return &r
}

func Opened(x int) Option {
	return func(c *Options) {
		c.Opened = x
	}
}

func Used(x compact.Instances) Option {
	return func(c *Options) {
		if c.Used == nil {
			c.Used = x
		} else {
			c.Used = c.Used.Merge(x)
		}
	}
}

func SetResults(x Results) Option {
	return func(c *Options) {
		c.Results = x
	}
}

func Declared(melds Melds) Option {
	opened := len(melds)
	return func(c *Options) {
		c.Declared = melds
		c.Opened += opened
	}
}

func Forms(x Form) Option {
	return func(c *Options) {
		c.Forms = x
	}
}

func StartMelds(m Melds) Option {
	return func(c *Options) {
		c.StartMelds = m
	}
}
