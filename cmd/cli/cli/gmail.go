package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/sjengpho/tin/grpc"
)

// NewGmailCommander returns a cli.GmailCommander.
func NewGmailCommander() GmailCommander {
	return &gmailCommander{}
}

// gmailCommander implements cli.GmailCommander.
type gmailCommander struct{}

// Login attempts to authorize the user.
//
// It generates an auth URL and asks the user to enter the authorization code.
func (s *gmailCommander) Login(c *grpc.Client) {
	authURL, err := c.GmailAuthURL()
	if err != nil {
		credentialsSuggestion := "Place the credentials.json at ~/.config/tin/gmail/credentials.json."
		if r, err := c.Config(); err == nil {
			fmt.Println(c)
			credentialsSuggestion = fmt.Sprintf("Place the credentials.json at %v", r.Config.GetGmailCredentials())
		}

		fmt.Println("Failed generating the authorization url")
		fmt.Println("")
		fmt.Println("- Enable the Gmail API (https://support.google.com/googleapi/answer/6158841?hl=en)")
		fmt.Printf("- %v\n", credentialsSuggestion)

		return
	}

	fmt.Printf("Visit the link below to retreive authorization code\n\n%v\n\n", authURL)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Authorization code: ")
	code, _ := reader.ReadString('\n')

	success := c.GmailAuthCode(code)
	if !success {
		log.Printf("failed authorizing using code: %v", code)
		return
	}
	fmt.Println("\nLogin success")
}

// Unread outputs the unread mail count.
func (s *gmailCommander) Unread(c *grpc.Client) {
	unread, err := c.GmailUnread()
	if err != nil {
		log.Printf("failed getting the unread mail count: %v", err)
		return
	}
	fmt.Println(unread)
}
