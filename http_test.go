package main

import (
	"testing"
)

func TestRenameFileName(t *testing.T) {
	tests := []struct {
		fileName string
		want     string
	}{
		{"random-string?test", "randomstringtest"},
		{"why\\/learn*golang?.mp4", "whylearngolang.mp4"},
		{"<2.2 - Why-do:tests/matter?>", "2.2  Whydotestsmatter"},
	}

	for _, v := range tests {
		if got := renameFileName(v.fileName); v.want != got {
			t.Errorf("renameFileName wanted: %s got: %s", v.want, got)
		}
	}
}
