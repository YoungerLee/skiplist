package skiplist

import (
	"testing"
)

func TestSkipList_Set(t *testing.T) {
	skiplist := NewSkipList[int, string](32, 4)
	skiplist.Set(1, "a")
	skiplist.Set(2, "b")
	skiplist.Set(3, "c")
	if skiplist.Length() != 3 {
		t.Errorf("expect length == 3, got %d", skiplist.Length())
		return
	}
	value := skiplist.Get(1)
	if value != "a" {
		t.Errorf("expect a, got %s", value)
		return
	}
	value = skiplist.Get(2)
	if value != "b" {
		t.Errorf("expect a, got %s", value)
		return
	}
	value = skiplist.Get(3)
	if value != "c" {
		t.Errorf("expect a, got %s", value)
		return
	}
	skiplist.Remove(3)
	if skiplist.Length() != 2 {
		t.Errorf("expect length == 2, got %d", skiplist.Length())
		return
	}
	value = skiplist.Get(3)
	if value != "" {
		t.Errorf("expect empty, got %s", value)
		return
	}
}
