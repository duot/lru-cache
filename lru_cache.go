package main

import (
	"errors"
	"fmt"
)

func main() {
	example := NewLRUCache(3)
	v, _ := example.get("aaa")
	fmt.Println(v, " should be -1")

	example.set("aaa", 123)
	v, _ = example.get("aaa")
	fmt.Println(v, " should be 123")

	example.remove("aaa")
	v, _ = example.get("aaa")
	fmt.Println(v, " should be -1")

	example.set("bbb", 12)
	example.set("ccc", 34)
	example.set("ddd", 56)
	//fmt.Println(example.list)
	fmt.Println(example.hash)

	example.set("eee", 78)
	//fmt.Println(example.list)
	fmt.Println(example.hash)

	example.set("fff", 99)
	fmt.Println(example.hash)

	example.remove("nonexistent")

	example.remove("fff")
	fmt.Println(example.hash)

	example.set("yyy", 0)
	fmt.Println(example.hash)

	example.set("zzz", 1)
	fmt.Println(example.hash)
}

type Lru struct {
	max  int
	hash map[string]*ListNode
	list *DLList
}

func NewLRUCache(maxSize int) *Lru {
	hash := map[string]*ListNode{}
	list := newDLList()
	cache := Lru{maxSize, hash, list}
	return &cache
}

func (cache *Lru) get(key string) (int, error) {
	node, found := cache.hash[key]
	if found {
		cache.setAsMostRecent(node)
		return node.data, nil
	}
	return -1, errors.New("no key")
}

func (cache *Lru) set(key string, data int) {
	node, ok := cache.hash[key]
	if ok {
		node.data = data
		cache.setAsMostRecent(node)
	} else {
		if cache.list.length == cache.max {
			cache.evictOldest()
		}
		node := cache.list.setNewestNode(key, data)
		cache.hash[key] = node
	}
}

func (cache *Lru) evictOldest() {
	delete(cache.hash, cache.list.oldest.key)
	cache.list.oldest = cache.list.oldest.prev
	cache.list.oldest.next = nil
	cache.list.length--
}

func (cache *Lru) setAsMostRecent(node *ListNode) {
	list := cache.list
	if node.prev != nil {
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	list.newest.prev = node
	list.newest = node
	node.next = list.newest
	node.prev = nil
}

func (cache *Lru) remove(key string) {
	node, found := cache.hash[key]
	if !found {
		return
	}

	// at front
	if node.prev == nil {
		newHead := cache.list.newest.next
		newHead.prev = nil
		cache.list.newest = newHead
	} else if node.next == nil {
		// at end
		newLast := node.prev
		newLast.next = nil
		cache.list.oldest = newLast
	} else {
		// somewhere in the middle
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	cache.list.length--
	delete(cache.hash, node.key)
}

type DLList struct {
	newest *ListNode
	oldest *ListNode
	length int
}

func newDLList() *DLList {
	list := DLList{nil, nil, 0}
	return &list
}

func (list *DLList) setNewestNode(key string, value int) *ListNode {
	node := &ListNode{key, value, list.newest, nil}
	if list.newest != nil {
		list.newest.prev = node
	}
	if list.length == 0 {
		list.oldest = node
	}
	list.newest = node
	list.length++
	return node
}

type ListNode struct {
	key  string
	data int
	next *ListNode
	prev *ListNode
}
