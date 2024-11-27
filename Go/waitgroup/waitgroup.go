//go:build !solution

package waitgroup

// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.

type Mutex struct {
	ch chan int
}

func (m *Mutex) Lock() {
	m.ch <- 0
}

func (m *Mutex) Unlock() {
	<-m.ch
}

func NewMutex() *Mutex {
	return &Mutex{make(chan int, 1)}
}

type WaitGroup struct {
	cnt   int
	ch    chan int
	Mutex *Mutex
}

// New creates WaitGroup.
func New() *WaitGroup {
	return &WaitGroup{0, make(chan int, 1), NewMutex()}
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically, this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
	wg.Mutex.Lock()
	defer wg.Mutex.Unlock()
	if wg.cnt == 0 {
		wg.ch <- 0
	}
	wg.cnt += delta
	if wg.cnt < 0 {
		panic("negative WaitGroup counter")
	}
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.Mutex.Lock()
	defer wg.Mutex.Unlock()
	wg.cnt--
	if wg.cnt < 0 {
		panic("negative WaitGroup counter")
	}
	if wg.cnt == 0 {
		<-wg.ch
	}
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	wg.Mutex.Lock()
	if wg.cnt == 0 {
		wg.Mutex.Unlock()
		return
	}
	wg.Mutex.Unlock()
	wg.ch <- 0
	<-wg.ch
}
