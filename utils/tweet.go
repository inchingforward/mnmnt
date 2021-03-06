package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/inchingforward/mnmnt/models"
)

// Tweet tweets out the title of the given memory along with a link
// to the memory's details page.  If any Twitter-related environment
// variables are not set, no tweet is attempted.
func Tweet(memory models.Memory) {
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

	body := fmt.Sprintf("%v %v/memories/%v", memory.Title, mnmntHost, memory.ID)

	tweet, err := api.PostTweet(body, nil)
	log.Printf("twitter result id: %v, error: %v\n", tweet.Id, err)
}
