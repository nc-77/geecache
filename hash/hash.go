package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Func func(data []byte) uint32

type Map struct {
	hash     Func
	replicas int            // 虚拟节点个数
	keys     []int          // 哈希环
	hashMap  map[int]string // 虚拟节点到物理节点的映射
}

func New(replicas int, hash Func) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hashKey := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hashKey)
			m.hashMap[hashKey] = key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	hashKey := int(m.hash([]byte(key)))
	n := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hashKey
	})
	if n == len(m.keys) {
		return m.hashMap[m.keys[0]]
	}
	return m.hashMap[m.keys[n]]
}
