package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {

	buffer := Init()

	assert.NoError(t, buffer.Update("counter", "a", 10))
	assert.NoError(t, buffer.Update("counter", "a", 10))

	m, ok := buffer.Get("counter", "a")
	assert.True(t, ok)
	assert.Equal(t, int64(20), *m.Delta)
	assert.Equal(t, "counter/a/20", m.String())

	assert.NoError(t, buffer.Update("gauge", "b", 10.0))
	m, ok = buffer.Get("gauge", "b")

	assert.True(t, ok)
	assert.Equal(t, float64(10.0), *m.Value)

	assert.NoError(t, buffer.Update("gauge", "b", 11.01))
	m, ok = buffer.Get("gauge", "b")

	assert.True(t, ok)
	assert.Equal(t, float64(11.01), *m.Value)
	assert.Equal(t, "gauge/b/11.01", m.String())

	err := buffer.Update("wrongtype", "a", 2000)
	assert.Error(t, err)

	m, ok = buffer.Get("gauge", "a")
	assert.Nil(t, m)
	assert.False(t, ok)
}

func TestBufferErrors(t *testing.T) {
	var err error

	buffer := Init()
	err = buffer.Update("nosuchtype", "a", 1000)
	assert.Error(t, err)

	err = buffer.Update("gauge", "b", "10")
	assert.Error(t, err)

	err = buffer.Update("counter", "c", 10.0)
	assert.Error(t, err)
}

func TestBufferFlush(t *testing.T) {
	buffer := Init()
	assert.NoError(t, buffer.Update("gauge", "b", 10.0))
	m, ok := buffer.Get("gauge", "b")

	assert.True(t, ok)
	assert.Equal(t, float64(10.0), *m.Value)

	buffer.Flush()
	m, ok = buffer.Get("gauge", "b")
	assert.Nil(t, m)
	assert.False(t, ok)
}
