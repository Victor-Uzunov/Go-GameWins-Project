package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	clientID     = "Ov23liuDotbDmTrCy3is"
	clientSecret = "806197c36f2b2c5e098bf9ee4d9a00358d1fb487"
	redirectURL  = "http://localhost:5000/callback"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	oauthURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email",
		clientID,
		url.QueryEscape(redirectURL),
	)

	http.Redirect(w, r, oauthURL, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	tokenURL := "https://github.com/login/oauth/access_token"
	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("code", code)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		http.Error(w, "Error decoding response", http.StatusInternalServerError)
		return
	}

	if tokenResponse.AccessToken == "" {
		http.Error(w, "No access token received: "+tokenResponse.Error, http.StatusInternalServerError)
		return
	}

	// Fetch user info
	userInfoURL := "https://api.github.com/user"
	req, err = http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Authorization", "Bearer "+tokenResponse.AccessToken)
	req.Header.Add("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Error fetching user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Error decoding user info", http.StatusInternalServerError)
		return
	}

	// Successfully fetched user info
	fmt.Fprintf(w, "User: %s (ID: %d)", userInfo.Login, userInfo.ID)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	fmt.Println("Server started at http://localhost:5000")
	http.ListenAndServe(":5000", nil)
}
