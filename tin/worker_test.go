package tin

import (
	"reflect"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	want := &Worker{}
	got := NewWorker(time.Second, func() {})

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestWorkerProcess(t *testing.T) {
	want := 7
	var got int

	worker := NewWorker(time.Millisecond, func() { got = 7 })
	time.Sleep(2 * time.Millisecond)

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}

	worker.Stop()
	time.Sleep(2 * time.Millisecond)
	want = 17
	got = 17
	time.Sleep(2 * time.Millisecond)

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
