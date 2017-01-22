package main

import (
	"github.com/petrkotek/go-google-contacts/contacts"
	"golang.org/x/oauth2"
	"fmt"
	"os"
	"context"
	"log"
)

var json_token_file = "token.json"

func NewAuthToken(){
	config := &oauth2.Config{
		ClientID: "369013585757-guuiqrkr32k2leb4kd2tspfff83g7dic.apps.googleusercontent.com",
		ClientSecret: "-TrVOp1NCd2a4voHv2BOLBAh",
		Endpoint: oauth2.Endpoint{
			AuthURL: "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
		Scopes: []string{"https://www.googleapis.com/auth/contacts.readonly"},
	}
	ctx := context.Background()
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	authStorage := &contacts.FileAuthStorage{json_token_file}
	authDetails := contacts.AuthDetails{
		AccessToken: tok.AccessToken,
		RefreshToken: tok.RefreshToken,
	}
	authStorage.Save(&authDetails)
}



func main() {
	if _, err := os.Stat(json_token_file); os.IsNotExist(err){
		NewAuthToken()
	}
	apiclient := contacts.NewClient(&contacts.StandardAuthManager{
		AuthStorage: &contacts.FileAuthStorage{json_token_file},
		AccessTokenRetriever: &contacts.StandardAccessTokenRetriever{
			ClientID: "",
			GoogleSecret: "",
		},
	})
	feed, err := apiclient.FetchFeed()
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
	for i, entry := range feed.Entries {
		fmt.Println("ENTRY ", i, ":")
		fmt.Printf("%+v\n\n", entry)
	}
}