package core

import "container/list"

// Element wraps a map key and value
type Element struct {
	key   string
	value *File
}

// LRUCacheMap is a fixed-size LRU hash map
type LRUCacheMap struct {
	hashMap map[string]*list.Element
	list    *list.List
	size    int
}

// NewLRUCacheMap returns a new LRUCacheMap of specified size
func NewLRUCacheMap(size int) *LRUCacheMap {
	return &LRUCacheMap{
		// size+1 to account for moment during put after adding new value but before old value is purged
		make(map[string]*list.Element, size+1),
		&list.List{},
		size,
	}
}

// Get returns file from LRUCacheMap for key
func (lru *LRUCacheMap) Get(key string) (*File, bool) {
	elem, ok := lru.hashMap[key]
	if !ok {
		return nil, ok
	}

	// Move element to front of the list
	lru.list.MoveToFront(elem)

	// Get Element and return *File value from it
	element, _ := elem.Value.(*Element)
	return element.value, ok
}

// Put file in LRUCacheMap at key
func (lru *LRUCacheMap) Put(key string, value *File) {
	elem := lru.list.PushFront(value)
	lru.hashMap[key] = elem

	if lru.list.Len() > lru.size {
		// Get element at back of list and Element from it
		elem = lru.list.Back()
		element, _ := elem.Value.(*Element)

		// Delete entry in hashMap with key from Element, and list
		delete(lru.hashMap, element.key)
		lru.list.Remove(elem)
	}
}

// Remove file in LRUCacheMap with key
func (lru *LRUCacheMap) Remove(key string) {
	elem, ok := lru.hashMap[key]
	if !ok {
		return
	}

	// Delete entry in hashMap with key from Element, and list
	element, _ := elem.Value.(*Element)
	delete(lru.hashMap, element.key)
	lru.list.Remove(elem)
}

// Iterate performsn an iteration over all key:value pairs in LRUCacheMap with supplied function
func (lru *LRUCacheMap) Iterate(iterator func(key string, value *File)) {
	for key := range lru.hashMap {
		element, _ := lru.hashMap[key].Value.(*Element)
		iterator(element.key, element.value)
	}
}
