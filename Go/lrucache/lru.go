//go:build !solution

package lrucache

import (
	"container/list"
)

func New(cap int) Cache {
	return &LRUCache{make(map[int]int, cap), cap, list.New()}
}

type LRUCache struct {
	m   map[int]int
	cap int
	l   *list.List
}

func (c *LRUCache) Get(key int) (int, bool) {
	v, ok := c.m[key]
	if ok {
		for e := c.l.Front(); e != nil; e = e.Next() {
			if e.Value == key {
				c.l.MoveToBack(e)
				break
			}
		}
	}
	return v, ok
}

func (c *LRUCache) Set(key int, value int) {
	_, ok := c.m[key]
	if ok {
		c.m[key] = value
		for e := c.l.Front(); e != nil; e = e.Next() {
			if e.Value == key {
				c.l.MoveToBack(e)
				break
			}
		}
	} else {
		if len(c.m) < c.cap {
			c.m[key] = value
			c.l.PushBack(key)
		} else if c.cap != 0 {
			removed := c.l.Remove(c.l.Front()).(int)
			delete(c.m, removed)
			c.m[key] = value
			c.l.PushBack(key)
		}
	}
}

func (c *LRUCache) Range(f func(key, value int) bool) {
	for e := c.l.Front(); e != nil; e = e.Next() {
		if !f(e.Value.(int), c.m[e.Value.(int)]) {
			break
		}
	}
}

func (c *LRUCache) Clear() {
	for k := range c.m {
		delete(c.m, k)
	}
	c.l = list.New()
}
