package tin

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// MailProvider is the interface implemented by an object that can
// return unread mails.
type MailProvider interface {
	UnreadMails() ([]Mail, error)
}

// Mail represents a mail message.
type Mail struct {
	Snippet string
}

// MailCount represents a mail count.
type MailCount int

// Equal implements tin.Comparable.
func (a MailCount) Equal(t interface{}) bool {
	if b, ok := t.(MailCount); ok {
		return a == b
	}
	return false
}

// UnreadMailCount represents a tin.StateKey.
const UnreadMailCount StateKey = "UnreadMailCount"

// MailService provides access to data from mail providers.
type MailService struct {
	provider MailProvider
	state    *State
	worker   *Worker
	logger   *log.Logger
}

// NewMailService returns a tin.MailService.
func NewMailService(p MailProvider, l *log.Logger) *MailService {
	s := &MailService{
		provider: p,
		state:    NewState(),
		logger:   l,
	}

	// Worker that fetches unread mails on intervals and updates the state.
	if p == nil {
		s.logger.Println(errors.New("failed initializing worker"))
	} else {
		s.worker = NewWorker(time.Minute, func() {
			mails, err := s.provider.UnreadMails()
			if err != nil {
				s.logger.Println(fmt.Errorf("worker failed: %w", err))
			} else {
				s.state.Set(UnreadMailCount, MailCount(len(mails)))
			}
		})
	}

	return s
}

// Subscribe returns a tin.StateSubscription.
func (s *MailService) Subscribe() StateSubscription {
	return s.state.Subscribe()
}

// UnreadMailCount returns a tin.MailCount.
func (s *MailService) UnreadMailCount() MailCount {
	v, err := s.state.Get(UnreadMailCount)
	if err != nil {
		return MailCount(0)
	}

	return v.(MailCount)
}

// SetUnreadMailCount updates the state.
func (s *MailService) SetUnreadMailCount(m MailCount) {
	s.state.Set(UnreadMailCount, m)
}
