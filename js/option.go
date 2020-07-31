package js

import "github.com/ssttevee/go-quickjs/internal"

type evalConfig struct {
	flags internal.EvalFlag
}

type EvalOption func(*evalConfig)

func evalOptionModule(c *evalConfig) {
	c.flags |= internal.EvalTypeModule
}

func EvalOptionStrict(c *evalConfig) {
	c.flags |= internal.EvalFlagStrict
}

func EvalOptionStrip(c *evalConfig) {
	c.flags |= internal.EvalFlagStrip
}

func EvalOptionBacktraceBarrier(c *evalConfig) {
	c.flags |= internal.EvalFlagBacktraceBarrier
}
