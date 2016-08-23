package snowflake

import (
	"testing"
)

// go test -bench="." -test.benchmem
func BenchmarkSnowFlake(t *testing.B) {
	uuid, _ := NewUUID(2)
	var pre int64 = 0

	for i := 0; i < t.N; i++ {
		id, err := uuid.Next()
		if err != nil {
			t.Log("ERROR:", err)
			continue
		}

		if id < 0 {
			t.Fatalf("id: %d < 0", id)
		}

		if pre >= id {
			t.Fatalf("pre: %d >= id: :%d", pre, id)
		}

		pre = id
	}
}
