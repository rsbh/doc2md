package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// GetConfig return oauth config
func GetConfig(clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

// GetClientCredentials return clientID and clientSecret
func GetClientCredentials() (string, string) {
	clientID := viper.GetString("CLIENT_ID")
	clientSecret := viper.GetString("CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("CLIENT_ID and CLIENT_SECRET is not available in env")
	}
	return clientID, clientSecret
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

// SaveToken save token to the path.
func SaveToken(path string, config *oauth2.Config) {
	token := getTokenFromWeb(config)
	saveToken(path, token)
	fmt.Println("Token saved at", path)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func tokenFromENV() (*oauth2.Token, error) {
	t := viper.GetString("TOKEN")
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(t), tok)
	if err != nil {
		return nil, err
	}
	return tok, err
}

// GetToken Get Token from env or file
func GetToken(tokenFile string) *oauth2.Token {
	var token *oauth2.Token
	if _, err := os.Stat(tokenFile); err == nil {
		token, _ = tokenFromFile(tokenFile)
	} else {
		token, _ = tokenFromENV()
	}
	if token == nil {
		log.Fatalln("Token is unavailable, Pass token from file or ENV")
	}
	return token
}
