package compiler

import (
	"sync"
)

type TaskRunner struct {
	work  chan func()
	group sync.WaitGroup
}

func (m *TaskRunner) Run(f func()) {
	m.group.Add(1)
	m.work <- f
}

func (m *TaskRunner) Kill() {
	close(m.work)
	m.Flush()
}

func (m *TaskRunner) Flush() {
	m.group.Wait()
}

func CreateTaskRunner(limit int) *TaskRunner {
	runner := &TaskRunner{
		work: make(chan func(), limit),
	}
	for i := 0; i < limit; i++ {
		go func() {
			for f := range runner.work {
				f()
				runner.group.Done()
			}
		}()
	}
	return runner
}
