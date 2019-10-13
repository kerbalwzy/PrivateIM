package utils

import "testing"

func TestBytesDataHash(t *testing.T) {
	strList := []string{"a", "ab", "abc"}
	for index, str := range strList {
		name := BytesDataHash([]byte(str))
		t.Logf("%d >> %s", index, name)
	}
}

