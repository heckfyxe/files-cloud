package transbytes

import "testing"

func TestSizeToIntAndString(t *testing.T) {
	assertEqual(0, 0, "B", t)
	assertEqual(1024, 1, "KB", t)
	assertEqual(1<<20, 1, "MB", t)
	assertEqual(1<<30, 1, "GB", t)
	assertEqual(1<<40, 1, "TB", t)
	assertEqual(1<<40+1<<30, 1.1, "TB", t)
	assertEqual(1<<40*10, 10, "TB", t)
}

func assertEqual(bytes float64, size float64, format string, t *testing.T) {
	if s, f := SizeToFloatAndString(bytes); s != size && f != format {
		t.Error("Excepted", size, format, ",but got", s, f)
	}
}
