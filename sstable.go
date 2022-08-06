package main

import "sync"

type SSTable struct {
	rwmu sync.RWMutex
	mt   *MemTable
}

func (s *SSTable) Set(k, v string) error {
	if !s.mt.IsFull() {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()
		return s.mt.Set(k, v)
	}
	return nil
	// number of keys
	// if key too much or size too big
	// lock
	// rename the current mt to a temp file (.temp)
	// create a new mt
	// unlock
	//
	// off load the mt to background
	// serialize to segment file
	// add segment file to managed segments (sorted by created time in descending order)
}
func (s *SSTable) Get() {
	// read from mt first
	// check segments' index file to find the key
	// return key if found
}
