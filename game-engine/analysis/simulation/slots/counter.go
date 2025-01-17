package magic

import (
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

type Counter interface {
	Name() string
	Done() bool
	Matches(res results.Results) bool
	Count() int
	Max() int
}

type Matcher func(res results.Results) bool

func NewCounter(name string, max int, matcher Matcher) Counter {
	return &counter{
		name:    name,
		max:     max,
		matcher: matcher,
	}
}

type counter struct {
	max     int
	count   int
	matcher Matcher
	name    string
	mutex   sync.RWMutex
}

func (c *counter) Name() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.name
}

func (c *counter) Done() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count >= c.max
}

func (c *counter) Matches(res results.Results) bool {
	matched := c.matcher(res)
	if matched {
		c.mutex.Lock()
		c.count++
		c.mutex.Unlock()
	}
	return matched
}

func (c *counter) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count
}

func (c *counter) Max() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.max
}
