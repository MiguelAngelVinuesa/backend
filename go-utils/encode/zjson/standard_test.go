package zjson

import (
	"encoding/json"
	"testing"
)

func BenchmarkStandardMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(d1)
	}
}

func BenchmarkStandardUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d2 := NewData1()
		json.Unmarshal(j1, d2)
		d2.ReturnToPool()
	}
}
