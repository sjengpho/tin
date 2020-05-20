package tin

import (
	"reflect"
	"testing"
)

type fakeNumber int

func (a fakeNumber) Equal(b interface{}) bool { return a == b }

func TestNewState(t *testing.T) {
	want := &State{}
	got := NewState()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestStateSetAndGet(t *testing.T) {
	s := NewState()
	s.Set("number", fakeNumber(7))
	s.Set("number", fakeNumber(7))

	tt := []struct {
		key  StateKey
		want StateValue
	}{
		{
			"number",
			fakeNumber(7),
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

	s.publish(StateMessage{key: key, value: fakeNumber(4)})
	s.publish(StateMessage{key: key, value: fakeNumber(17)})
	s.publish(StateMessage{key: key, value: fakeNumber(1)})
	s.publish(StateMessage{key: key, value: fakeNumber(2)})
	s.publish(StateMessage{key: key, value: fakeNumber(3)})
	s.publish(StateMessage{key: key, value: fakeNumber(4)})
	s.publish(StateMessage{key: key, value: fakeNumber(7)})

	want := fakeNumber(7)
	got := <-subscription.Channel
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestStateSubscriptionClose(t *testing.T) {
	s := NewState()
	a := s.Subscribe()
	s.Subscribe()
	a.Close()

	want := 1
	got := len(s.subscribers)
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
