package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func main() {
	ctx := context.Background()

	credentials, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secrets file: %v", err)
	}

	config, err := google.ConfigFromJSON(credentials, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secrets to config: %v", err)
	}

	client := getClient(config)

	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Gmail client: %v", err)
	}

	user := "me"
	r, err := gmailService.Users.Messages.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	fmt.Println("Messages:")
	if len(r.Messages) == 0 {
		fmt.Println("No messages found.")
	} else {
		for _, m := range r.Messages {
			fmt.Printf("- Message ID: %s\n", m.Id)
		}
	}
}
