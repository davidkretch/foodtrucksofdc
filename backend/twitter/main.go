package main

import (
	"fmt"
	"log"
)

func main() {
	key := ""
	secret := ""

	token, err := getToken(key, secret)
	if err != nil {
		log.Fatal(err)
	}
	tweets, err := getTweets("davidkretch", 0, token)
	if err != nil {
		log.Fatal(err)
	}

	for _, tweet := range tweets {
		fmt.Println(tweet.Text)
	}
}
