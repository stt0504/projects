//go:build !solution

package tparallel

import "sync"

var wgSeq sync.WaitGroup
var wgParallel sync.WaitGroup

type T struct {
	ch         chan struct{}
	parallelCh chan struct{}
	wg         sync.WaitGroup
	main       bool
	mainCh     chan struct{}
	parent     *T
	subtests   []*T
	wgParallel sync.WaitGroup
}

func (t *T) Parallel() {
	if t.parent != nil {
		t.parent.wgParallel.Add(1)
	}
	if t.parent != nil {
		t.parent.subtests = append(t.parent.subtests, t)
	}
	t.parallelCh = make(chan struct{})
	if t.ch != nil {
		close(t.ch)
	}
	if t.main {
		wgSeq.Done()
		wgParallel.Add(1)
	}
	if t.parent != nil {
		t.parent.wg.Done()
	}
	t.parallelCh <- struct{}{}
	if t.parent != nil {
		t.parent.wgParallel.Done()
	}
}

func (t *T) Run(subtest func(t *T)) {
	t.wg.Add(1)
	subT := &T{mainCh: make(chan struct{}), subtests: make([]*T, 0), parent: t}

	go func(subtest func(t *T)) {
		subtest(subT)
		if subT.parallelCh == nil {
			t.wg.Done()
		}
	}(subtest)

	t.wg.Wait()

	for _, subtest := range subT.subtests {
		<-subtest.parallelCh
	}
	subT.wgParallel.Wait()
}

func Run(topTests []func(t *T)) {
	ch := make([]chan struct{}, len(topTests))
	var tests []*T

	for i := range topTests {
		ch[i] = make(chan struct{}, 1)
	}

	for i, test := range topTests {
		wgSeq.Add(1)
		t := &T{
			ch:       ch[i],
			main:     true,
			mainCh:   make(chan struct{}),
			subtests: make([]*T, 0, 10),
		}
		tests = append(tests, t)

		go func(test func(t *T), i int) {
			if i > 0 {
				<-ch[i-1]
			}
			test(t)
			if t.parallelCh == nil {
				wgSeq.Done()
				close(t.ch)
			} else {
				wgParallel.Done()
			}
			for _, subtest := range t.subtests {
				<-subtest.parallelCh
			}
			t.wgParallel.Wait()
		}(test, i)
	}

	wgSeq.Wait()
	for _, t := range tests {
		if t.parallelCh != nil {
			<-t.parallelCh
		}
	}

	wgParallel.Wait()
}
