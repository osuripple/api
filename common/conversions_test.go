package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafeByteStringConversions(t *testing.T) {
	type tc struct {
		name string
		b    []byte
		s    string
	}
	for _, cas := range []tc{
		{"simple", []byte("hello"), "hello"},
		{"new lines", []byte("hello\nworld"), "hello\nworld"},
	} {
		t.Run(cas.name, func(t *testing.T) {
			t.Run("b2s", func(t *testing.T) {
				v := b2s(cas.b)
				assert.Equal(t, cas.s, v)
				assert.Equal(t, len(cas.s), len(v))
			})

			t.Run("s2b", func(t *testing.T) {
				v := s2b(cas.s)
				assert.Equal(t, cas.b, v)
				assert.Equal(t, len(cas.b), len(v))
			})
		})
	}
}
