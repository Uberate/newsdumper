package staged

import (
	"context"
	"fmt"
	"sync"
)

func NewConcurrentStage(name string, ignoreError bool, stageds ...Staged) Staged {
	return NewConcurrentStageWithKindVersion("default-concurrent-stage", "v1", name, ignoreError, stageds...)
}

func NewConcurrentStageWithKindVersion(kind, version, name string, ignoreError bool, stageds ...Staged) Staged {
	return &ConcurrentStaged{
		kind:              kind,
		version:           version,
		concurrentStageds: stageds,
		ignoreError:       ignoreError,
		name:              name,
	}
}

type ConcurrentStaged struct {
	kind              string
	version           string
	name              string
	concurrentStageds []Staged
	nextStaged        Staged
	ignoreError       bool
}

func (c *ConcurrentStaged) Kind() string {
	return c.kind
}

func (c *ConcurrentStaged) Name() string {
	return c.name
}

func (c *ConcurrentStaged) Version() string {
	return c.version
}

func (c *ConcurrentStaged) RunE(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(c.concurrentStageds))

	errors := []error{}

	for _, item := range c.concurrentStageds {
		item := item
		go func() {
			if err := item.RunE(ctx); err != nil {
				errors = append(errors, fmt.Errorf("staged: kind: [%s], version: %s, name: [%s] has err: %v",
					item.Kind(), item.Version(), item.Name(), err))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if !c.ignoreError && len(errors) != 0 {
		return fmt.Errorf("errors: %v", errors)
	}

	if c.nextStaged != nil {
		return c.nextStaged.RunE(ctx)
	}

	return nil
}

func (c *ConcurrentStaged) SetNext(staged Staged) {
	c.nextStaged = staged
}
