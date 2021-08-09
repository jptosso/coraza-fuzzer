package fuzzer

import "testing"

func TestContentSize(t *testing.T) {
	size := uint64(50)
	for f, fn := range Rules {
		data := fn(size)
		if len(data) < int(size) {
			t.Errorf("%s size does not match, got %d", f, len(data))
		}
	}
}
