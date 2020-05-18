package tin

import (
	"reflect"
	"testing"
	"time"
)

func TestNewState(t *testing.T) {
	want := &State{}
	got := NewState()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestStateSetAndGet(t *testing.T) {
	s := NewState()
	s.Set("number", 7)

	tt := []struct {
		key  StateKey
		want StateValue
	}{
		{
			"number",
			7,
		},
		{
			"invalid",
			nil,
		},
	}

	for _, tc := range tt {
		got, err := s.Get(tc.key)

		if err != nil && got != nil {
			t.Errorf("want %v, got %v", tc.want, got)
		}

		if got != tc.want {
			t.Errorf("want %v, got %v", tc.want, got)
		}
	}
}

func TestStateSubscribe(t *testing.T) {
	s := NewState()
	want := make(<-chan StateValue, 1)
	got := s.Subscribe()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestStatePublish(t *testing.T) {
	s := NewState()
	ch := s.Subscribe()
	s.publish(4)
	s.publish(17)
	s.publish(7)

	time.Sleep(time.Millisecond * 100)

	want := 7
	got := <-ch
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
