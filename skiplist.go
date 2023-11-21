package skiplist

import (
	"golang.org/x/exp/constraints"
	"math/rand"
	"sync"
)

type SkipNode[K constraints.Ordered, V any] struct {
	Key   K
	Value V
	Next  []*SkipNode[K, V] // [0, level)
}

type SkipList[K constraints.Ordered, V any] struct {
	SkipNode[K, V] // head
	mu             sync.RWMutex
	maxLevel       int
	skip           int
	level          int
	length         int
	updates        []*SkipNode[K, V]
}

func NewSkipList[K constraints.Ordered, V any](maxLevel, skip int) *SkipList[K, V] {
	return &SkipList[K, V]{
		maxLevel: maxLevel,
		skip:     skip,
		SkipNode: SkipNode[K, V]{
			Next: make([]*SkipNode[K, V], maxLevel),
		},
		updates: make([]*SkipNode[K, V], maxLevel),
	}
}

func (l *SkipList[K, V]) Get(key K) (value V) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	curr := &l.SkipNode
	for i := l.level - 1; i >= 0; i-- {
		for curr != nil && curr.Key < key {
			curr = curr.Next[i]
		}
	}
	if curr != nil && curr.Key == key {
		value = curr.Value
	}
	return
}

func (l *SkipList[K, V]) Set(key K, value V) {
	l.mu.Lock()
	defer l.mu.Unlock()
	prev := &l.SkipNode
	var curr *SkipNode[K, V]
	for i := l.level - 1; i >= 0; i-- {
		curr = prev.Next[i]
		for curr != nil && curr.Key < key {
			prev = curr
			curr = curr.Next[i]
		}
		l.updates[i] = prev
	}
	if curr != nil && curr.Key == key {
		curr.Value = value
		return
	}
	if l.level < l.maxLevel && (l.level == 0 || toSkip(l.skip)) {
		l.level = min(l.level+1, l.maxLevel)
		l.updates[l.level-1] = &l.SkipNode
	}
	node := &SkipNode[K, V]{
		Key:   key,
		Value: value,
		Next:  make([]*SkipNode[K, V], l.level),
	}
	for i := 0; i < l.level; i++ {
		node.Next[i] = l.updates[i].Next[i]
		l.updates[i].Next[i] = node
	}
	l.length++
}

func (l *SkipList[K, V]) Remove(key K) (value V) {
	l.mu.Lock()
	defer l.mu.Unlock()
	prev := &l.SkipNode
	var curr *SkipNode[K, V]
	for i := l.level - 1; i >= 0; i-- {
		curr = prev.Next[i]
		for curr != nil && curr.Key < key {
			prev = curr
			curr = curr.Next[i]
		}
		l.updates[i] = prev
	}
	if curr == nil || curr.Key != key {
		return
	}
	for i, node := range curr.Next {
		if l.updates[i].Next[i] == curr {
			l.updates[i].Next[i] = node
			if l.SkipNode.Next[i] == nil {
				l.level--
			}
		}
	}
	l.length--
	l.updates = make([]*SkipNode[K, V], l.level)
	return
}

func (l *SkipList[K, V]) Length() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.length
}

func toSkip(skip int) bool {
	return rand.Int()%skip == 0
}
