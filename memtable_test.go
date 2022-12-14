package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MemTable(t *testing.T) {
	mt, err := NewMemTable(
		WithNamePrefix(t.Name()),
	)
	require.NoError(t, err)

	require.NoError(t, mt.Set("a", "1"))
	val, err := mt.Get("a")
	require.NoError(t, err)
	require.Equal(t, "1", val)

	require.NoError(t, mt.Set("a", "2"))
	val, err = mt.Get("a")
	require.NoError(t, err)
	require.Equal(t, "2", val)

	val, err = mt.Get("b")
	require.Error(t, err)
	require.Empty(t, val)

	require.NoError(t, mt.Set("b", "123"))
	require.NoError(t, mt.Set("c", "456"))
	require.NoError(t, mt.Close())

	// reopen
	mt, err = NewMemTable(
		WithNamePrefix(t.Name()),
	)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"a": "2",
		"b": "123",
		"c": "456",
	}, mt.m)
}

func Test_MemTable_IsFull(t *testing.T) {
	mt, err := NewMemTable(
		WithNamePrefix(t.Name()),
		WithKeyThreshold(3),
	)
	require.NoError(t, err)
	require.False(t, mt.IsFull())

	require.NoError(t, mt.Set("a", "1"))
	require.False(t, mt.IsFull())

	require.NoError(t, mt.Set("b", "2"))
	require.False(t, mt.IsFull())

	require.NoError(t, mt.Set("c", "3"))
	require.True(t, mt.IsFull())
	require.NoError(t, mt.Close())

	mt, err = NewMemTable(
		WithNamePrefix(t.Name()),
		WithSizeThreshold(20), // 20 bytes
	)
	require.NoError(t, err)
	require.False(t, mt.IsFull())

	require.NoError(t, mt.Set("d", strings.Repeat("1", 20)))
	require.True(t, mt.IsFull(), "current size: %v", mt.size)
}

func Test_MemTable_Rename(t *testing.T) {
	mt, err := NewMemTable(
		WithNamePrefix(t.Name()),
	)
	require.NoError(t, err)
	require.NoError(t, mt.Set("a", "1"))
	require.NoError(t, mt.Set("b", "2"))
	require.NoError(t, mt.Set("c", "3"))
	require.NoError(t, mt.Rename())
	require.NoError(t, mt.Close())

	// reopen again without loading temp file
	mt, err = NewMemTable(
		WithNamePrefix(t.Name()),
		WithUseTempFile(false),
	)
	require.NoError(t, err)
	require.Equal(t, map[string]string{}, mt.m)
	require.NoError(t, mt.Set("c", "4")) // will overwrite "c" to "3"
	require.Equal(t, map[string]string{"c": "4"}, mt.m)
	require.NoError(t, mt.Close())

	// reopen again with loading temp file
	mt, err = NewMemTable(
		WithNamePrefix(t.Name()),
		WithUseTempFile(true),
	)
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"a": "1",
		"b": "2",
		"c": "4", // "c" should not equal to "3"
	}, mt.m)
}
