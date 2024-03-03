package util

import "errors"

// https://forum.golangbridge.org/t/im-searching-for-something-like-a-both-way-map/2557/8

// A BidirectionalMap is a map-like data structure that can be
// searched either by key or by value.
type BidirectionalMap struct {
	m1 map[string]int
	m2 map[int]string
}

// Insert inserts a key/value pair into the BidirectionalMap.
func (m *BidirectionalMap) Insert(key string, value int) {
	m.m1[key], m.m2[value] = value, key
}

// FindByKey searches for a string in a BidirectionalMap based on map keys.
func (m *BidirectionalMap) FindByKey(key string) (int, error) {
	if val, found := m.m1[key]; found {
		return val, nil
	}
	return 0, errors.New("key not found")
}

// FindByValue searches for a string in a BidirectionalMap based on map values.
func (m *BidirectionalMap) FindByValue(value int) (string, error) {
	if key, found := m.m2[value]; found {
		return key, nil
	}
	return "", errors.New("value not found")
}
