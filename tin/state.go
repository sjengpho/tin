package tin

import (
	"errors"
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
	values map[StateKey]StateValue

	// Pubsub.
	subscribers map[chan StateValue]chan StateValue
	msgCh       chan StateMessage
}

// StateKey represents the key of a state value.
type StateKey string

// StateValue represents the value of a state.
type StateValue interface {
	Comparable
}

// StateMessage holds a StateKey and StateValue.
type StateMessage struct {
	key   StateKey
	value StateValue
}

// StateSubscription contains a read-only channel that receives changes of every state value.
//
// Close removes the subscription and closes the channel.
type StateSubscription struct {
	Channel <-chan StateValue
	Close   func()
}

// NewState returns tin.State
func NewState() *State {
	s := &State{
		values:      make(map[StateKey]StateValue),
		subscribers: make(map[chan StateValue]chan StateValue),
		msgCh:       make(chan StateMessage, 1),
	}

	go s.work()

	return s
}

// Get returns a tin.StateValue for the given tin.StateKey.
//
// An error will be returned if an entry doesn't exists.
func (s *State) Get(k StateKey) (StateValue, error) {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.values[k]
	if !ok {
		return nil, ErrEntryNotExist
	}

	return v, nil
}

// Set updates the state and sends the value to the subscribers.
func (s *State) Set(k StateKey, v StateValue) {
	s.Lock()
	if old, exists := s.values[k]; !exists || !v.Equal(old) {
		s.values[k] = v
		s.publish(StateMessage{key: k, value: v})
	}
	s.Unlock()
}

// Subscribe creates and returns a tin.StateSubscription.
func (s *State) Subscribe() StateSubscription {
	s.Lock()
	ch := make(chan StateValue, 1)
	s.subscribers[ch] = ch
	s.Unlock()

	return StateSubscription{
		Channel: ch,
		Close: func() {
			s.Lock()
			close(ch)
			delete(s.subscribers, ch)
			s.Unlock()
		},
	}
}

// publish sends the message to the channel.
func (s *State) publish(m StateMessage) {
	s.msgCh <- m
}

// work proccesses messages and sends them to subscribers after a given duration.
//
// Messages of the same type will reset the duration and replace the previous message
// which effectively mean in case of bursts, only the last message will be sent to the
// subscribers.
func (s *State) work() {
	pending := struct {
		sync.Mutex
		messages map[StateKey]*time.Timer
	}{
		messages: map[StateKey]*time.Timer{},
	}

	for msg := range s.msgCh {
		pending.Lock()
		if m, exists := pending.messages[msg.key]; exists {
			m.Stop() // Stopping the timer from sending the message.
		}

		// Adding a pending message, potentially replacing the previous one.
		pending.messages[msg.key] = func(msg StateMessage) *time.Timer {
			return time.AfterFunc(time.Second, func() {
				s.RLock()
				for _, s := range s.subscribers {
					s <- msg.value
				}
				s.RUnlock()

				pending.Lock()
				delete(pending.messages, msg.key) // Cleanup after sending the message.
				pending.Unlock()
			})
		}(msg)
		pending.Unlock()
	}
}
