package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sjengpho/tin/tin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var readFile = ioutil.ReadFile

// Service implements tin.UnreadMailCounter.
type Service struct {
	tokenPath  string
	configPath string
	client     *http.Client
}

// NewService returns a new Service.
func NewService(c string, t string) *Service {
	return &Service{
		configPath: c,
		tokenPath:  t,
	}
}

// UnreadMails fetches unread messages from Gmail.
//
// It assumes setting "is:unread" will filter unread mails only.
// Because of the pagination it will keep requesting messages
// until NextPageToken is empty.
func (s *Service) UnreadMails() ([]tin.Mail, error) {
	client, err := s.getUsersMessagesService()
	if client == nil {
		return []tin.Mail{}, err
	}

	resp, err := client.List("me").Q("is:unread").Do()
	if err != nil {
		return []tin.Mail{}, fmt.Errorf("Failed fetching unread mails: %w", err)
	}

	msgs := resp.Messages
	for resp.NextPageToken != "" {
		resp, err = client.List("me").Q("is:unread").PageToken(resp.NextPageToken).Do()
		msgs = append(msgs, resp.Messages...)
	}

	mm := []tin.Mail{}
	for _, msg := range msgs {
		mm = append(mm, tin.Mail{
			Snippet: msg.Snippet,
		})
	}

	return mm, nil
}

// AuthURL uses oauth2.Config to return a URL to the consent page.
func (s *Service) AuthURL() (string, error) {
	config, err := s.getConfig()
	if err != nil {
		return "", err
	}

	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), nil
}

// ExchangeAuthCode attempts to exchange the authorization code.
//
// The token will be stored as JSON for future uses.
func (s *Service) ExchangeAuthCode(code string) error {
	config, err := s.getConfig()
	if err != nil {
		return err
	}

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return fmt.Errorf("failed getting the token: %w", err)
	}

	if err := s.saveToken(token); err != nil {
		return fmt.Errorf("failed saving the token: %w", err)
	}

	return nil
}

// getUsersMessagesService returns gmail.UsersMessagesService.
func (s *Service) getUsersMessagesService() (*gmail.UsersMessagesService, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return gmail.NewUsersMessagesService(service), nil
}

// getClient creates and returns a http.Client.
//
// The client will be cached for future uses.
func (s *Service) getClient() (*http.Client, error) {
	if s.client != nil {
		return s.client, nil
	}

	config, err := s.getConfig()
	if err != nil {
		return nil, err
	}

	token, err := s.getToken()
	if err != nil {
		return nil, err
	}

	s.client = config.Client(context.Background(), token)

	return s.client, nil
}

// getConfig attempts to read and parse the config file into a oauth2.Config.
func (s *Service) getConfig() (*oauth2.Config, error) {
	bytes, err := readFile(s.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed opening file %v: %w", s.configPath, err)
	}

	config, err := google.ConfigFromJSON(bytes, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("failed parsing credentials file: %w", err)
	}

	return config, nil
}

// getToken attempts to read and parse the token file into a oauth2.Token.
//
// It assumes that the token is JSON encoded.
func (s *Service) getToken() (*oauth2.Token, error) {
	bytes, err := readFile(s.tokenPath)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	err = json.Unmarshal(bytes, token)

	return token, err
}

// saveToken encodes the token into a JSON and stores it at the tokenPath.
func (s *Service) saveToken(t *oauth2.Token) error {
	f, err := os.OpenFile(s.tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed saving the token: %w", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(t)

	return nil
}
