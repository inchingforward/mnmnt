package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/inchingforward/mnmnt/models"
)

func Tweet(memory *models.Memory) {
	mnmntHost := os.Getenv("MONUMENT_HOST")
	consumerKey := os.Getenv("MONUMENT_TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("MONUMENT_TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("MONUMENT_TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("MONUMENT_TWITTER_ACCESS_SECRET")

	if mnmntHost == "" || consumerKey == "" || consumerSecret == "" || accessToken == "" || accessTokenSecret == "" {
		log.Println("Missing mail environment variables...not tweeting")
		return
	}

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	body := fmt.Sprintf("%v %v/memories/%v", memory.Title, mnmntHost, memory.Id)

	tweet, err := api.PostTweet(body, nil)
	log.Printf("twitter result id: %v, error: %v\n", tweet.Id, err)
}
