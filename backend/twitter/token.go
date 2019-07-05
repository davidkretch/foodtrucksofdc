package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// getToken gets a Twitter API bearer token.
// See https://developer.twitter.com/en/docs/basics/authentication/overview/application-only
func getToken(key string, secret string) (string, error) {
	credentials := getTokenCredentials(key, secret)
	req := buildTokenRequest(credentials)
	token, err := sendTokenRequest(req)
	if err != nil {
		return "", err
	}
	return token, nil
}

// getTokenCredentials translates consumer key and secret into bearer
// token credentials.
func getTokenCredentials(key string, secret string) string {
	k := url.QueryEscape(key)
	s := url.QueryEscape(secret)
	return base64.StdEncoding.EncodeToString([]byte(k + ":" + s))
}

// buildTokenRequest builds the bearer token HTTP request.
func buildTokenRequest(credentials string) *http.Request {
	reqBody := strings.NewReader("grant_type=client_credentials")
	req, _ := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", reqBody)
	req.Header.Add("Authorization", "Basic "+credentials)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	return req
}

// sendTokenRequest sends the request and returns the token.
func sendTokenRequest(req *http.Request) (string, error) {
	resp, err := sendRequest(req)
	if err != nil {
		return "", err
	}
	message, err := decodeTokenResponse(resp)
	if err != nil {
		return "", err
	}
	return message, nil
}

// A tokenResponse stores the Twitter API's token-request response.
type tokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

// decodeTokenResponse unmarshals the response body.
func decodeTokenResponse(resp []byte) (string, error) {
	var msg tokenResponse
	err := json.Unmarshal(resp, &msg)
	if err != nil {
		return "", err
	}
	if msg.TokenType != "bearer" {
		return "", fmt.Errorf("did not receive bearer token. received %s token", msg.TokenType)
	}
	return msg.AccessToken, nil
}
