package tin

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

type mailProviderMock struct {
	returnError bool
}

func (m mailProviderMock) UnreadMails() ([]Mail, error) {
	if m.returnError {
		return nil, errors.New("error")
	}

	return make([]Mail, 0), nil
}

func TestNewMailService(t *testing.T) {
	tt := []struct {
		want *MailService
		got  *MailService
	}{
		{
			want: &MailService{},
			got:  NewMailService(mailProviderMock{returnError: false}, log.New(os.Stdout, "", log.Flags())),
		},
		{
			want: &MailService{},
			got:  NewMailService(mailProviderMock{returnError: true}, log.New(os.Stdout, "", log.Flags())),
		},
		{
			want: &MailService{},
			got:  NewMailService(nil, log.New(os.Stdout, "", log.Flags())),
		},
	}

	for _, tc := range tt {
		if reflect.TypeOf(tc.got) != reflect.TypeOf(tc.want) {
			t.Errorf("want %v, got %v", reflect.TypeOf(tc.want), reflect.TypeOf(tc.got))
		}
	}
}

func TestMailSubscribe(t *testing.T) {
	s := NewMailService(mailProviderMock{}, log.New(os.Stdout, "", log.Flags()))
	want := StateSubscription{}
	got := s.Subscribe()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestMailUnreadMailCount(t *testing.T) {
	withState := NewMailService(nil, log.New(ioutil.Discard, "", log.Flags()))
	withState.SetUnreadMailCount(MailCount(0))

	tt := []struct {
		service *MailService
		want    MailCount
	}{
		{
			service: withState,
			want:    MailCount(0),
		},
		{
			service: NewMailService(nil, log.New(ioutil.Discard, "", log.Flags())),
			want:    MailCount(0),
		},
	}

	for _, tc := range tt {
		got := tc.service.UnreadMailCount()

		if got != tc.want {
			t.Errorf("want %v, got %v", tc.want, got)
		}
	}
}

func TestMailCountEqualTrue(t *testing.T) {
	a := MailCount(1)
	b := MailCount(1)

	want := true
	got := a.Equal(b)
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestMailCountEqualFalse(t *testing.T) {
	tt := []struct {
		a MailCount
		b interface{}
	}{
		{
			a: MailCount(1),
			b: MailCount(2),
		},
		{
			a: MailCount(1),
			b: 1,
		},
		{
			a: MailCount(1),
			b: 2,
		},
	}

	for _, tc := range tt {
		want := false
		got := tc.a.Equal(tc.b)

		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}
