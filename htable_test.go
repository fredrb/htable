package htable

import (
	"testing"
)

const (
	Key1       = StringKey("key1")
	Key2       = StringKey("key2")
	MissingKey = StringKey("missing key")
	NumberKey  = IntKey(31)
)

var ht Table

func assertTableWithKey(t *testing.T, key PreHashable) interface{} {
	value, ok := ht.Get(key)
	if !ok {
		t.Errorf("Key1 couldn't be retrieved with .Get")
	}
	return value
}

func assertTableWithKeyAndStringValue(t *testing.T, key PreHashable, expected string) {
	v, ok := assertTableWithKey(t, key).(string)
	if !ok {
		t.Errorf("Failed to convert value to string")
	}
	if v != expected {
		t.Errorf("Expected value to be %s but instead found %s", expected, v)
	}
}

func assertLength(t *testing.T, expected int) {
	if length := ht.Len(); length != expected {
		t.Errorf("Expected size to be %d but instead found %d", expected, ht.Len())
	}
}

func TestCrateTable(t *testing.T) {
	ht = New()
	item, _ := ht.Get(MissingKey)
	if item != nil {
		t.Errorf("Shouldn't have any item with key, but item found")
	}
}

func TestAddKeyToTable(t *testing.T) {
	ht = New()
	ht.Set(Key1, "This is my value")
	assertTableWithKeyAndStringValue(t, Key1, "This is my value")
}

func TestOverwriteKey(t *testing.T) {
	ht = New()
	ht.Set(Key1, "First value")
	ht.Set(Key1, "Second value")
	assertTableWithKeyAndStringValue(t, Key1, "Second value")
}

func TestAddMultipleKeysOfDifferentTypes(t *testing.T) {
	ht = New()
	ht.Set(Key1, "Value1")
	ht.Set(NumberKey, "Value2")
	assertTableWithKeyAndStringValue(t, Key1, "Value1")
	assertTableWithKeyAndStringValue(t, NumberKey, "Value2")
}

func TestDeleteValueFromHash(t *testing.T) {
	ht = New()
	ht.Set(Key1, "v1")
	err := ht.Delete(Key1)
	if err != nil {
		t.Errorf("Shouldn't have failed when deleting from hash")
	}
	_, ok := ht.Get(Key1)
	if ok {
		t.Errorf("Should have deleted Key1")
	}
}

func TestDeleteErrorIfKeyIsMissing(t *testing.T) {
	ht = New()
	err := ht.Delete(MissingKey)
	if err == nil {
		t.Errorf("Should have returned error")
	}
}

func TestTableSize(t *testing.T) {
	ht = New()
	ht.Set(Key1, "v1")
	assertLength(t, 1)
	ht.Set(Key2, "v2")
	assertLength(t, 2)
	_ = ht.Delete(Key1)
	assertLength(t, 1)
}

func TestReadAnyValueFromInterface(t *testing.T) {
	ht = New()
	type point struct {
		x int
		y int
	}

	ht.Set(Key1, point{1, 2})
	ht.Set(Key2, point{3, 4})

	p, _ := ht.Get(Key1)
	value, ok := p.(point)
	if !ok {
		t.Errorf("Couldn't convert point")
	}
	if value.x != 1 || value.y != 2 {
		t.Errorf("Point converted but values don't match")
	}

	p, _ = ht.Get(Key2)
	value, ok = p.(point)
	if !ok {
		t.Errorf("Couldn't convert point")
	}
	if value.x != 3 || value.y != 4 {
		t.Errorf("Point converted but values don't match")
	}
}

// This test only makes sense if threshold for resizing is 0.5 and size doubles every resize
func TestAutoIncreaseSize(t *testing.T) {
	ht = NewSized(2)
	ht.Set(Key1, "v1")
	ht.Set(Key2, "v2")
	ht.Set(NumberKey, "v3")

	if len(ht.buckets) < 8 {
		t.Errorf("Should have doubled Table size twice: current size: %d", len(ht.buckets))
	}

	ht.Set(StringKey("abc1"), 1)
	ht.Set(StringKey("abc2"), 2)
	ht.Set(StringKey("abc3"), 3)
	ht.Set(StringKey("abc4"), 4)
	ht.Set(StringKey("abc5"), 5)
	ht.Set(StringKey("abc6"), 6)

	if len(ht.buckets) < 16 {
		t.Errorf("Should have doubled Table size up to 16: current size: %d", len(ht.buckets))
	}
}
