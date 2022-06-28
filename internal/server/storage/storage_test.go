package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageDumpToFile(t *testing.T) {
	assert.Equal(t, true, true)

	st := Init()

	st.Update("gauge", "test_gauge", 69.420)
	st.Update("counter", "test_counter", int64(42))

	require.NoError(t, st.DumpToFile("/tmp/ololo.txt"))

	st2, err := InitFromFile("/tmp/ololo.txt")
	require.NotNil(t, st2)
	require.NoError(t, err)

	m1, ok := st2.Get("gauge", "test_gauge")
	require.True(t, ok)
	m2, ok := st2.Get("counter", "test_counter")
	require.True(t, ok)

	assert.Equal(t, 69.420, *m1.Value)
	assert.Equal(t, int64(42), *m2.Delta)

	os.Remove("/tmp/ololo.txt")
}
