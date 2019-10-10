package livescore

import (
	"context"
	"sync"
)

type Worker struct {
	url string
	mu sync.Mutex
	champs []Champ
	err error
	ctx context.Context
	cancel context.CancelFunc
	wg sync.WaitGroup
}

type Workers struct {
	xs map[string] *Worker
	mu sync.Mutex
	ctx context.Context
	cancel context.CancelFunc
}

func NewWorkers(ctx context.Context) *Workers{
	x := &Workers{
		xs: make(map[string]*Worker),
	}
	x.ctx, x.cancel = context.WithCancel(ctx)
	return x
}

func (x *Workers) Close() {
	x.cancel()
	x.mu.Lock()
	defer x.mu.Unlock()
	for _,w := range x.xs{
		w.Cancel()
	}
}

func (x *Workers) Get(url string) *Worker{
	x.mu.Lock()
	defer x.mu.Unlock()
	if w,f := x.xs[url]; f {
		return w
	}
	w := newWorker(url,x.ctx)
	x.xs[url] = w
	return w
}

func newWorker(url string, ctx context.Context) *Worker{
	x := &Worker{ url: url,}
	x.ctx, x.cancel = context.WithCancel(ctx)
	x.wg.Add(1)
	x.mu.Lock()

	go func() {
		defer x.wg.Done()

		x.err = FetchChamps(x.url, &x.champs)
		x.mu.Unlock()

		log.Info("worker started: " + x.url)
		for{
			select {
			case <-ctx.Done():
				return
			default:
				var champs []Champ
				err := FetchChamps(x.url, &champs)
				x.mu.Lock()
				x.err = err
				x.champs = append([]Champ{}, champs...)
				x.mu.Unlock()
			}
		}
	}()
	return x
}

func (x *Worker) Cancel()  {
	x.cancel()
	x.wg.Wait()
}

func (x *Worker) Champs() ([]Champ,error)  {
	x.mu.Lock()
	defer x.mu.Unlock()
	return append([]Champ{}, x.champs...), x.err
}
