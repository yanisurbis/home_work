package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Run("random stuff", func(t *testing.T) {
		assert.Equal(t, 123, 123, "they should be equal")
	})
}

