package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	_memTableFilenamePrefix = ".memtable"
	_LogFile                = ".log"
	_LogLine                = "%s|%v\n"
	_SegmentFile            = ".%s.segment"
	_IndexFile              = ".%s.index"

	_unlimited            uint = 0
	_defaultSizeThreshold uint = 4 * 1024 * 1024 * 1024 // 4 MiB
)

type memTableOptions struct {
	namePrefix    string
	keyThreshold  uint
	sizeThreshold uint
}

type Option interface {
	apply(m *memTableOptions)
}

type namePrefixOption string

func (o namePrefixOption) apply(m *memTableOptions) { m.namePrefix = string(o) }
func WithNamePrefix(n string) Option                { return namePrefixOption(n) }

type keyThresholdOption uint

func (o keyThresholdOption) apply(m *memTableOptions) { m.keyThreshold = uint(o) }
func WithKeyThreshold(k uint) Option                  { return keyThresholdOption(k) }

type sizeThresholdOption uint

func (o sizeThresholdOption) apply(m *memTableOptions) { m.sizeThreshold = uint(o) }
func WithSizeThreshold(s uint) Option                  { return sizeThresholdOption(s) }

type MemTable struct {
	f             *os.File
	m             map[string]string
	keyThreshold  uint
	sizeThreshold uint
	size          uint
}

func NewMemTable(opt ...Option) (*MemTable, error) {
	options := memTableOptions{
		namePrefix:    _memTableFilenamePrefix,
		keyThreshold:  _unlimited,
		sizeThreshold: _unlimited,
	}
	for _, opt := range opt {
		opt.apply(&options)
	}
	if options.keyThreshold == _unlimited && options.sizeThreshold == _unlimited {
		options.sizeThreshold = _defaultSizeThreshold
	}

	m, n, err := rebuildMemTable(options.namePrefix + _LogFile)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(options.namePrefix+_LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	mt := &MemTable{
		m:             m,
		f:             f,
		keyThreshold:  options.keyThreshold,
		sizeThreshold: options.sizeThreshold,
		size:          n,
	}
	return mt, nil
}

func (m MemTable) Close() error {
	return m.f.Close()
}

func (m *MemTable) Set(k, v string) error {
	n, err := m.f.Write([]byte(fmt.Sprintf(_LogLine, k, v)))
	if err != nil {
		return nil
	}

	m.size += uint(n)
	m.m[k] = v
	return nil
}

func (m MemTable) Get(k string) (string, error) {
	v, ok := m.m[k]
	if !ok {
		return "", fmt.Errorf("key %s not found", k)
	}
	return v, nil
}

func (m MemTable) IsFull() bool {
	return (m.keyThreshold != _unlimited && uint(len(m.m)) >= m.keyThreshold) ||
		(m.sizeThreshold != _unlimited && m.size >= m.sizeThreshold)
}

// func (m MemTable) Serialize() ([]byte, []byte) {
// 	keys := make([]string, 0, len(m.m))
// 	for k := range m.m {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)

// 	offset := 0
// 	index := make([]byte, 0)
// 	segment := make([]byte, 0)
// 	for _, k := range keys {
// 		index = append(index, []byte(fmt.Sprintf("%s|%v\n", k, offset))...)
// 		pair := []byte(fmt.Sprintf("%s|%v\n", k, m.m[k]))
// 		offset += len(pair)
// 		segment = append(segment, pair...)
// 	}
// 	return index, segment
// }

// func (m MemTable) Write(filename string) error {
// 	idxf, err := os.OpenFile(fmt.Sprintf(_IndexFile, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		return err
// 	}
// 	segf, err := os.OpenFile(fmt.Sprintf(_SegmentFile, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		return err
// 	}
// 	index, segment := m.Serialize()
// 	if _, err := idxf.Write(index); err != nil {
// 		return err
// 	}
// 	if _, err := segf.Write(segment); err != nil {
// 		return err
// 	}
// 	return nil
// }

func rebuildMemTable(filename string) (map[string]string, uint, error) {
	m := make(map[string]string)

	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) /* something went wrong */ {
		return nil, 0, err
	} else if os.IsNotExist(err) {
		// file does not exist, return empty map
		return m, 0, nil
	}

	// then, means file exists
	rfd, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}

	// rebuild the memtable from the log file.
	var n uint
	scanner := bufio.NewScanner(rfd)
	for scanner.Scan() {
		b := scanner.Bytes()
		n += uint(len(b)) + 1 // add new line character
		parts := strings.Split(string(b), "|")
		m[parts[0]] = parts[1]
	}
	return m, n, nil
}
