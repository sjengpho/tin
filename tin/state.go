package tin

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

// ErrInvalidType means the type of the new value is not equal to
// the type of the initial value.

// ErrEntryNotExist means an entry doesn't exists for the given key.
var ErrEntryNotExist = errors.New("Entry doesn't exist")

// State represents the state.
//
// Subscribers get state updates on state changes.
// In case of bursts only the last message will be sent to the subscribers.
type State struct {
	sync.RWMutex
	values map[StateKey]interface{}

	// Pubsub.
	subscribers []chan interface{}
	msgCh       chan interface{}
}

// StateKey represents the key of an entry of state values.
type StateKey string

// NewState returns tin.State
func NewState() *State {
	s := &State{
		values:      make(map[StateKey]interface{}),
		subscribers: []chan interface{}{},
		msgCh:       make(chan interface{}, 1),
	}

	go s.work()

	return s
}

// Get returns the state value for the given tin.StateKey.
//
// An error will be returned if an entry doesn't exists.
func (s *State) Get(k StateKey) (interface{}, error) {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.values[k]
	if !ok {
		return nil, ErrEntryNotExist
	}

	return v, nil
}

// Set updates the state and sends the value to the subscribers.
func (s *State) Set(k StateKey, v interface{}) {
	s.Lock()
	s.values[k] = v
	s.Unlock()
	s.publish(v)
}

// Subscribe creates and return a read-only channel that can receive state updates.
func (s *State) Subscribe() <-chan interface{} {
	s.Lock()
	ch := make(chan interface{}, 1)
	s.subscribers = append(s.subscribers, ch)
	s.Unlock()

	return ch
}

// publish sends the message to the channel.
func (s *State) publish(msg interface{}) {
	s.msgCh <- msg
}

// work proccesses messages and sends them to subscribers after a given duration.
//
// Messages of the same type will reset the duration and replace the previous message
// which effectively mean in case of bursts, only the last message will be sent to the
// subscribers.
func (s *State) work() {
	pending := struct {
		sync.Mutex
		messages map[string]*time.Timer
	}{
		messages: map[string]*time.Timer{},
	}

	for v := range s.msgCh {
		id := reflect.TypeOf(v).String()

		pending.Lock()
		if m, exists := pending.messages[id]; exists {
			m.Stop() // Stopping the timer from sending the message.
		}

		// Adding a pending message, potentially replacing the previous one.
		pending.messages[id] = func(msg interface{}) *time.Timer {
			return time.AfterFunc(time.Second, func() {
				for _, s := range s.subscribers {
					s <- msg
				}

				pending.Lock()
				delete(pending.messages, id) // Cleanup after sending the message.
				pending.Unlock()
			})
		}(v)
		pending.Unlock()
	}
}
