package tin

import (
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	want := Config{}
	got := DefaultConfig()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}
