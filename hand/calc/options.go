package calc

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/meld"
)

type Options struct {
	Opened  int
	Used    compact.Instances
	Melds   meld.Melds
	Results Results
}

type Option func(*Options)

func GetOptions(opts ...Option) *Options {
	var r Options
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

func Melds(melds meld.Melds) Option {
	used := compact.NewInstances()
	melds.AddTo(used)
	opened := len(melds)
	return func(c *Options) {
		c.Melds = melds
		c.Opened = opened
		Used(used)(c)
	}
}
