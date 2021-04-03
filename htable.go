package htable

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	initialM                   = 16
	fnvOffSetBasis      uint64 = 14695981039346656037
	fnvPrime                   = 1099511628211
	loadFactorThreshold        = 0.5
)

type IntKey int
type StringKey string

type PreHashable interface {
	HashBytes() []byte
	Equal(PreHashable) bool
}

func (i IntKey) HashBytes() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, int64(i))
	return buf[:n]
}

func (i IntKey) Equal(other PreHashable) bool {
	v, ok := other.(IntKey)
	return ok && i == v
}

func (str StringKey) HashBytes() []byte {
	return []byte(str)
}

func (str StringKey) Equal(other PreHashable) bool {
	v, ok := other.(StringKey)
	return ok && str == v
}

type Table struct {
	length  int
	buckets [][]entry
}

type entry struct {
	key   PreHashable
	value interface{}
}

// Create a new Hash Table hinting the desired number of buckets
func NewSized(initial int) Table {
	return Table{
		buckets: make([][]entry, initial),
	}
}

// Create a new Hash Table
func New() Table {
	return NewSized(initialM)
}

func hashValue(v PreHashable, limit int) int {
	hash := fnvOffSetBasis
	for _, b := range v.HashBytes() {
		hash = hash ^ uint64(b)
		hash = hash * fnvPrime
	}
	return int(hash % uint64(limit))
}

func (ht *Table) expandTable() error {
	newTable := make([][]entry, len(ht.buckets)*2)
	for _, bucket := range ht.buckets {
		for _, e := range bucket {
			newHash := hashValue(e.key, len(ht.buckets))
			newTable[newHash] = append(newTable[newHash], entry{e.key, e.value})
		}
	}
	ht.buckets = newTable
	return nil
}

func (ht *Table) loadFactor() float32 {
	return float32(ht.length) / float32(len(ht.buckets))
}

func (ht *Table) Set(key PreHashable, value interface{}) {
	hash := hashValue(key, len(ht.buckets))
	// check if key is already added, if yes, just overwrite
	for i, e := range ht.buckets[hash] {
		if e.key == key {
			ht.buckets[hash][i].value = value
			return
		}
	}

	ht.buckets[hash] = append(ht.buckets[hash], entry{key, value})
	ht.length += 1
	if ht.loadFactor() > loadFactorThreshold {
		ht.expandTable()
	}
}

func (ht *Table) Get(key PreHashable) (interface{}, bool) {
	hash := hashValue(key, len(ht.buckets))
	for _, v := range ht.buckets[hash] {
		if v.key == key {
			return v.value, true
		}
	}
	return nil, false
}

func (ht *Table) Len() int {
	return ht.length
}

func (ht *Table) Delete(key PreHashable) error {
	hash := hashValue(key, len(ht.buckets))
	for i, v := range ht.buckets[hash] {
		if v.key == key {
			current := ht.buckets[hash]
			current[i] = current[len(current)-1]
			current = current[:len(current)-1]
			ht.length -= 1
			ht.buckets[hash] = current
			return nil
		}
	}

	return fmt.Errorf("Key error")
}

func (ht *Table) Dump(w io.Writer) {
	fmt.Fprintf(w, "length = %d\n", ht.length)
	for i, entries := range ht.buckets {
		fmt.Fprintf(w, "bucket %3d: ", i)
		for j, entry := range entries {
			fmt.Fprintf(w, "%s:%v", entry.key, entry.value)
			if j < len(entries)-1 {
				fmt.Fprintf(w, ", ")
			}
		}
		fmt.Fprintln(w)
	}
}
