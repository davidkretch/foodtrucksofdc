package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// getTweets gets tweets for the given user.
// See https://developer.twitter.com/en/docs/tweets/timelines/api-reference/get-statuses-user_timeline.html.
func getTweets(user string, since uint64, token string) ([]Tweet, error) {
	req := buildTweetRequest(user, since, token)
	resp, err := sendTweetRequest(req)
	if err != nil {
		return []Tweet{}, nil
	}
	return resp, nil
}

func buildTweetRequest(user string, since uint64, token string) *http.Request {
	u, _ := url.Parse("https://api.twitter.com/1.1/statuses/user_timeline.json")
	q := u.Query()
	q.Set("screen_name", user)
	// q.Set("since_id", strconv.FormatUint(since, 10))
	u.RawQuery = q.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Authorization", "Bearer "+token)
	return req
}

func sendTweetRequest(req *http.Request) ([]Tweet, error) {
	resp, err := sendRequest(req)
	if err != nil {
		return nil, err
	}
	tweets, err := decodeTweetResponse(resp)
	if err != nil {
		return []Tweet{}, nil
	}
	return tweets, nil
}

// A Tweet contains information about a tweet, e.g. text and timestamp.
type Tweet struct {
	CreatedAt string `json:"created_at"`
	ID        uint64 `json:"id"`
	Text      string `json:"text"`
}

func decodeTweetResponse(resp []byte) ([]Tweet, error) {
	var tweets []Tweet
	err := json.Unmarshal(resp, &tweets)
	if err != nil {
		return []Tweet{}, err
	}
	return tweets, nil
}
