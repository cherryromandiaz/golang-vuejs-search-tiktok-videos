package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

// Video represents the video details we want to extract
type Video struct {
	URL      string `json:"url"`
	Age      string `json:"age"`
	//ImageURL string `json:"imageUrl"`
	
	//VideoSrc string `json:"videoSrc"`
}

func main() {
	headless := flag.Bool("headless", true, "Run in headless mode")
	flag.Parse()

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// Allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Query parameter is missing", http.StatusBadRequest)
			return
		}

		// Fetch videos from TikTok based on the query
		videos, err := fetchTikTokVideos(query, *headless)
		if err != nil {
			http.Error(w, "Failed to fetch videos", http.StatusInternalServerError)
			log.Printf("Error fetching videos: %v", err)
			return
		}

		// Send results as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(videos)
	})

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// fetchTikTokVideos uses Chromedp to perform a TikTok search and scrape video data
func fetchTikTokVideos(query string, headless bool) ([]Video, error) {
	searchURL := fmt.Sprintf("https://www.tiktok.com/search?q=%s", query)
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	}
	if headless {
		opts = append(opts, chromedp.Headless)
	}

	// Create a new context with Chromedp
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Define the slice to hold our results
	var videos []Video

	// Set a timeout to prevent long-running operations
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Run Chromedp tasks to navigate to the TikTok search page and scrape data
	err := chromedp.Run(ctx,
		chromedp.Navigate(searchURL),
		chromedp.WaitVisible(`[data-e2e="search_top-item"]`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Use JS within Chromedp to get video data
			return chromedp.Evaluate(`(() => {
				const items = Array.from(document.querySelectorAll('[data-e2e="search_top-item"]'));
				return items.map(item => {
					const link = item.querySelector('a') ? item.querySelector('a').href : '';
					//const img = item.querySelector('img') ? item.querySelector('img').src : '';
					const age = item.querySelector('.css-dennn6-DivTimeTag') ? item.querySelector('.css-dennn6-DivTimeTag').innerText : '';
					//const videoSrcElem = item.querySelector('.css-1fofj7p-DivBasicPlayerWrapper .xgplayer-container video');
					//const videoSrc = videoSrcElem ? videoSrcElem.src : '';
					//return { url: link, imageUrl: img, age: age, videoSrc: videoSrc };
					return { url: link, age: age };
				});
			})()`, &videos).Do(ctx)
		}),
	)

	if err != nil {
		return nil, err
	}

	return videos, nil
}
