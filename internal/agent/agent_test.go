package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	agent := Init(2*time.Second, 3*time.Second, nil, nil)
	assert.NotNil(t, agent)
}
