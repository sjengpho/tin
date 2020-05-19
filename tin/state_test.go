package tin

import (
	"reflect"
	"testing"
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
	want := StateSubscription{}
	got := s.Subscribe()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestStatePublish(t *testing.T) {
	s := NewState()
	subscription := s.Subscribe()
	key := StateKey("key")

	s.publish(StateMessage{key: key, value: 4})
	s.publish(StateMessage{key: key, value: 17})
	s.publish(StateMessage{key: key, value: 1})
	s.publish(StateMessage{key: key, value: 2})
	s.publish(StateMessage{key: key, value: 3})
	s.publish(StateMessage{key: key, value: 4})
	s.publish(StateMessage{key: key, value: 7})

	want := 7
	got := <-subscription.Channel
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
