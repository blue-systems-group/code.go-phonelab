package phonelab

import (
	"reflect"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	sorted_loglines, err := Parse("fixtures/sorted.out")
	if err != nil {
		t.Error(err)
	}
	if len(sorted_loglines) != 11 {
		t.Errorf("Logline length mismatch: %v != 11", len(sorted_loglines))
	}
	unsorted_loglines, err := Parse("fixtures/unsorted.out")
	if err != nil {
		t.Error(err)
	}
	if len(unsorted_loglines) != 11 {
		t.Errorf("Logline length mismatch: %v != 11", len(unsorted_loglines))
	}
	if reflect.DeepEqual(sorted_loglines, unsorted_loglines) {
		t.Error("Unsorted loglines match sorted loglines.")
	}
	sort.Sort(ByTime(unsorted_loglines))
	if !reflect.DeepEqual(sorted_loglines, unsorted_loglines) {
		t.Error("Sorted loglines do not match.")
	}
	sort.Sort(ByTime(sorted_loglines))
	if !reflect.DeepEqual(sorted_loglines, unsorted_loglines) {
		t.Error("Sorted loglines do not match after unnecessary sort.")
	}
}
