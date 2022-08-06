package main

import "fmt"

func must(msg string, err error) {
	if err != nil {
		panic("failed to " + msg + ": " + err.Error())
	}
}
func main() {
	fmt.Println("hello world")
	t, err := NewMemTable()
	must("create memtable", err)
	t.Set("a", "1")
	t.Set("b", "2")
}
