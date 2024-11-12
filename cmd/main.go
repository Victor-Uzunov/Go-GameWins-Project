package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	oauth2Config *oauth2.Config
	oauth2State  = "random" // This should be a random string in production
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	// Set up OAuth2 configuration
	oauth2Config = &oauth2.Config{
		ClientID:     getGithubClientID(),
		ClientSecret: getGithubClientSecret(),
		RedirectURL:  "http://localhost:4000/login/github/callback",
		Scopes:       []string{"read:org", "read:user"},
		Endpoint:     github.Endpoint,
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login/github/", githubLoginHandler)
	http.HandleFunc("/login/github/callback", githubCallbackHandler)
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		loggedinHandler(w, r, "")
	})

	fmt.Println("[ UP ON PORT 4000 ]")
	log.Panic(http.ListenAndServe(":4000", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/login/github/">LOGIN</a>`)
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauth2State)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Panic("Token exchange failed", err)
	}

	accessToken := token.AccessToken
	githubData := getGithubData(accessToken)

	loggedinHandler(w, r, githubData)
}

func getGithubClientID() string {
	clientID, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}
	fmt.Println(clientID)
	return clientID
}

func getGithubClientSecret() string {
	clientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	if !exists {
		log.Fatal("Github Client Secret not defined in .env file")
	}
	fmt.Println(clientSecret)
	return clientSecret
}

func getGithubData(accessToken string) string {
	userData := getGithubUserData(accessToken)
	userOrgs := getGithubUserOrgs(accessToken)
	role := determineUserRole(userOrgs)

	fullData := fmt.Sprintf(`{
		"user": %s,
		"role": "%s",
		"organizations": %s
	}`, userData, role, userOrgs)

	return fullData
}

func getGithubUserData(accessToken string) string {
	client := oauth2Config.Client(oauth2.NoContext, &oauth2.Token{AccessToken: accessToken})
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Panic("Request failed", err)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)
	return string(respbody)
}

func getGithubUserOrgs(accessToken string) string {
	client := oauth2Config.Client(oauth2.NoContext, &oauth2.Token{AccessToken: accessToken})
	resp, err := client.Get("https://api.github.com/user/orgs")
	if err != nil {
		log.Panic("Request failed", err)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)
	return string(respbody)
}

func determineUserRole(orgs string) string {
	var orgList []struct {
		Login string `json:"login"`
	}
	json.Unmarshal([]byte(orgs), &orgList)

	for _, org := range orgList {
		if org.Login == "Reader-Role" {
			return "reader"
		} else if org.Login == "Writer-Role" {
			return "writer"
		} else if org.Login == "Admin-Role" {
			return "admin"
		}
	}

	return "guest"
}

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "application/json")

	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}
