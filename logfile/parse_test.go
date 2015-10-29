package phonelab

import (
	"testing"
)

func TestParse(t *testing.T) {
	loglines, err := Parse("fixtures/20-2000.out.gz")
	if err != nil {
		t.Error(err)
	}
	if len(loglines) != 2000 {
		t.Errorf("Logline length mismatch: %v != 2000", len(loglines))
	}
	loglines, err = Parse("fixtures/20-1000.out")
	if err != nil {
		t.Error(err)
	}
	if len(loglines) != 1000 {
		t.Errorf("Logline length mismatch: %v != 1000", len(loglines))
	}
}
