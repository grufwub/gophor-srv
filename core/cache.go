package core

import "container/list"

// element wraps a map key and value
type element struct {
	key   string
	value *file
}

// lruCacheMap is a fixed-size LRU hash map
type lruCacheMap struct {
	hashMap map[string]*list.Element
	list    *list.List
	size    int
}

// newLRUCacheMap returns a new LRUCacheMap of specified size
func newLRUCacheMap(size int) *lruCacheMap {
	return &lruCacheMap{
		// size+1 to account for moment during put after adding new value but before old value is purged
		make(map[string]*list.Element, size+1),
		&list.List{},
		size,
	}
}

// Get returns file from LRUCacheMap for key
func (lru *lruCacheMap) Get(key string) (*file, bool) {
	lElem, ok := lru.hashMap[key]
	if !ok {
		return nil, ok
	}

	// Move element to front of the list
	lru.list.MoveToFront(lElem)

	// Get Element and return *File value from it
	element, _ := lElem.Value.(*element)
	return element.value, ok
}

// Put file in LRUCacheMap at key
func (lru *lruCacheMap) Put(key string, value *file) {
	lElem := lru.list.PushFront(&element{key, value})
	lru.hashMap[key] = lElem

	if lru.list.Len() > lru.size {
		// Get element at back of list and Element from it
		lElem = lru.list.Back()
		element, _ := lElem.Value.(*element)

		// Delete entry in hashMap with key from Element, and from list
		delete(lru.hashMap, element.key)
		lru.list.Remove(lElem)
	}
}

// Remove file in LRUCacheMap with key
func (lru *lruCacheMap) Remove(key string) {
	lElem, ok := lru.hashMap[key]
	if !ok {
		return
	}

	// Delete entry in hashMap and list
	delete(lru.hashMap, key)
	lru.list.Remove(lElem)
}

// Iterate performs an iteration over all key:value pairs in LRUCacheMap with supplied function
func (lru *lruCacheMap) Iterate(iterator func(key string, value *file)) {
	for key := range lru.hashMap {
		element, _ := lru.hashMap[key].Value.(*element)
		iterator(element.key, element.value)
	}
}
