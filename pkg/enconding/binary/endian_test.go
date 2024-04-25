package binary

import "testing"

func TestInt64(t *testing.T) {
	buf := make([]byte, 4)
	BigEndian.PutInt32(buf, 100)
	i := BigEndian.Int32(buf)
	t.Logf("i: %d", i)
}
