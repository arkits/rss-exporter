package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/arkits/rss-api/models"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/viper"
)

type SavedFeedData struct {
	lastUpdated time.Time
	pollCount   int
	savedFeed   *gofeed.Feed
}

type FeedsStore struct {
	mu             sync.Mutex
	savedFeedsData map[string]SavedFeedData
}

var feedStore FeedsStore

// BeginPollingFeeds kicks of individual goroutines to poll the RSS feeds
func BeginPollingFeeds() {

	// Get the feed URLs from config
	feedURLs := viper.GetStringSlice("rss.feeds")
	log.Printf("Retrived feeds to poll from config - %v", feedURLs)

	// Initialize the savedFeedsData Map
	sfd := make(map[string]SavedFeedData)

	// Update the feedStore with initialized savedFeedsData Map
	feedStore.mu.Lock()
	feedStore.savedFeedsData = sfd
	feedStore.mu.Unlock()

	// Iterate through the feedURLs and spawn goroutines to begin the polling
	for _, feedURL := range feedURLs {
		go pollFeed(feedURL)
	}
}

// Polls the feedURL and updated the feedStore
func pollFeed(feedURL string) {
	for {

		timeStart := time.Now()

		parsedFeed := parseFeed(feedURL)

		if parsedFeed != nil {
			// Update the SavedFeedData
			sfd := feedStore.savedFeedsData[feedURL]
			sfd.lastUpdated = time.Now()
			sfd.pollCount++
			sfd.savedFeed = parsedFeed

			// Update the feedStore
			feedStore.mu.Lock()
			feedStore.savedFeedsData[feedURL] = sfd
			feedStore.mu.Unlock()

			log.Printf("Updated SavedFeedData for feedURL=%v", feedURL)
		}

		feedPollDuration := fmt.Sprintf(`feed_poll_duration{url="%v"}`, feedURL)
		metrics.GetOrCreateSummary(feedPollDuration).UpdateDuration(timeStart)

		time.Sleep(viper.GetDuration("rss.pollRateMs") * time.Millisecond)
	}
}

func parseFeed(feedURL string) *gofeed.Feed {

	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(feedURL)

	if err != nil {
		log.Printf("Caught Error when parsing feed! feedURL=%v error=%v", feedURL, err.Error())
		return nil
	}

	return parsedFeed
}

// FeedHandler returns the parsedFeed based on the URL
func FeedHandler(w http.ResponseWriter, r *http.Request) {

	feedURL := r.URL.Query().Get("url")
	if feedURL == "" {
		var error models.Error
		error.Message = "Param url not provided!"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}

	log.Printf("feedURL=%v", feedURL)

	parsedFeed := feedStore.savedFeedsData[feedURL].savedFeed
	if parsedFeed == nil {
		var error models.Error
		error.Message = "Feed not available!"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsedFeed)
	return
}
