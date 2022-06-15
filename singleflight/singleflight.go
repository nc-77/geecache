package singleflight

import "sync"

type call struct {
	wg    sync.WaitGroup
	value interface{}
	err   error
}

type Group struct {
	mu    sync.Mutex
	calls map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.calls == nil {
		g.calls = make(map[string]*call)
	}
	if c, ok := g.calls[key]; ok {
		g.mu.Unlock()
		// 正在运行call，其他协程等待结果
		c.wg.Wait()
		return c.value, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.mu.Unlock()

	c.value, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()

	return c.value, c.err
}
