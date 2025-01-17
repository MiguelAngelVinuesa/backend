package zjson

import (
	"testing"

	"github.com/goccy/go-json"
)

func BenchmarkGoccyMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(d1)
	}
}

func BenchmarkGoccyUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d2 := NewData1()
		json.Unmarshal(j1, d2)
		d2.ReturnToPool()
	}
}
