package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	_LogFile     = ".memtable.log"
	_LogLine     = "%s|%v\n"
	_SegmentFile = ".memtable.%s.segment"
	_IndexFile   = ".memtable.%s.index"
)

// Write --->
// Open, restore

// TODO: add log file for avoid data loss when app crashed.
type MemTable struct {
	m map[string]string
	f *os.File
}

func parseLogLine(line string) (string, string) {
	parts := strings.Split(line, "|")
	return parts[0], parts[1]
}

func rebuildMemTable(filename string) (map[string]string, error) {
	m := make(map[string]string)

	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) /* something went wrong */ {
		return nil, err
	} else if os.IsNotExist(err) {
		// file does not exist, return empty map
		return m, nil
	}

	// then, means file exists
	rfd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// rebuild the memtable from the log file.
	scanner := bufio.NewScanner(rfd)
	for scanner.Scan() {
		line := scanner.Text()
		k, v := parseLogLine(line)
		m[k] = v
	}

	return m, nil
}

func NewMemTable() (*MemTable, error) {
	m, err := rebuildMemTable(_LogFile)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(_LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	mt := &MemTable{
		m: m,
		f: f,
	}
	return mt, nil
}

func (m MemTable) Set(k, v string) error {
	_, err := m.f.Write([]byte(fmt.Sprintf(_LogLine, k, v)))
	if err != nil {
		return nil
	}

	m.m[k] = v
	return nil
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
