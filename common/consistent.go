package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type uints []uint32

func (x uints) Len() int {
	return len(x)
}

func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("hash 环没有数据")

type Consistent struct {
	circle map[uint32]string
	sortedHashes uints
	VirtualNode int
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		circle:       make(map[uint32]string),
		VirtualNode:  20,
	}
}

func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}


func (c *Consistent) hashkey(key string) uint32 {
	if len(key) < 64 {
		var srcatch [64]byte
		copy(srcatch[:], key)
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	if cap(c.sortedHashes) / (c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}

	for k := range c.circle {
		hashes = append(hashes, k)
	}

	sort.Sort(hashes)

	c.sortedHashes = hashes
}

func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()

	c.add(element)
}

func (c *Consistent) add(element string) {
	for i:=0;i<c.VirtualNode;i++ {
		c.circle[c.hashkey(c.generateKey(element, i))] = element
	}

	c.updateSortedHashes()
}

func (c *Consistent) remove (element string) {
	for i:=0;i<c.VirtualNode;i++ {
		delete(c.circle, c.hashkey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()

	c.remove(element)
}

func (c *Consistent) search(key uint32) int {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}

	//使用二分查找算法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}

	return 0
}

func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return "", errEmpty
	}

	key := c.hashkey(name)
	i := c.search(key)

	return c.circle[c.sortedHashes[i]], nil
}