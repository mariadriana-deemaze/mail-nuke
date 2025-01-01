package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {
	const tokenFile = "token.json"
	token, err := readTokenFromFile(tokenFile)
	if err != nil {
		token = getTokenFromWeb(config)
		saveTokenToFile(tokenFile, token)
	}
	return config.Client(context.Background(), token)
}

func readTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	return &token, err
}

func saveTokenToFile(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to save token to file: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL to authenticate: \n%v\n", authURL)

	var authCode string
	fmt.Print("Enter the authorization code: ")
	fmt.Scan(&authCode)

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

func extractSpammySenders() []string {
	file, err := os.Open("mails.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var emails []string
	if err := json.Unmarshal(data, &emails); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return emails
}

func deleteEmailsBySenders(srv *gmail.Service, senders []string) {
	user := "me"
	batchSize := 1000 // Gmail API only allows up to 1000 message IDs per call

	for _, sender := range senders {
		query := fmt.Sprintf("from:%s", sender)
		fmt.Printf("Searching for emails from: %s\n", sender)

		var pageToken string
		var messageIDs []string

		for {
			call := srv.Users.Messages.List(user).Q(query).MaxResults(500)
			if pageToken != "" {
				call = call.PageToken(pageToken)
			}

			messages, err := call.Do()
			if err != nil {
				log.Printf("Unable to retrieve messages for sender %s: %v", sender, err)
				break
			}

			for _, msg := range messages.Messages {
				messageIDs = append(messageIDs, msg.Id)
				if len(messageIDs) >= batchSize {
					err := batchDeleteMessages(srv, user, messageIDs)
					if err != nil {
						log.Printf("Error deleting messages for sender %s: %v", sender, err)
					}
					messageIDs = []string{}
				}
			}

			pageToken = messages.NextPageToken
			if pageToken == "" {
				break
			}
		}

		if len(messageIDs) > 0 {
			err := batchDeleteMessages(srv, user, messageIDs)
			if err != nil {
				log.Printf("Error deleting remaining messages for sender %s: %v", sender, err)
			}
		}

		fmt.Printf("Completed processing emails from: %s\n", sender)
	}
}

func batchDeleteMessages(srv *gmail.Service, user string, messageIDs []string) error {
	req := &gmail.BatchDeleteMessagesRequest{
		Ids: messageIDs,
	}
	err := srv.Users.Messages.BatchDelete(user, req).Do()
	if err != nil {
		return fmt.Errorf("batch delete failed: %w", err)
	}
	fmt.Printf("Batch deleted %d messages\n", len(messageIDs))
	return nil
}

func main() {
	ctx := context.Background()

	credentials, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secrets file: %v", err)
	}

	config, err := google.ConfigFromJSON(credentials, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secrets to config: %v", err)
	}

	client := getClient(config)

	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Gmail client: %v", err)
	}

	spammySenders := extractSpammySenders()

	deleteEmailsBySenders(gmailService, spammySenders)
}
