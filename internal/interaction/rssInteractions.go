package interaction

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := getFeedRequest(ctx, feedURL)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	res, err := getFeedResponse(req)
	if err != nil {
		fmt.Println("Error getting response:", err)
		os.Exit(1)
	}
	data, err := readFeedResponse(res)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}
	rss := RSSFeed{}

	if err := xml.Unmarshal(data, &rss); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		os.Exit(1)
	}
	xmlUnescape(&rss)

	return &rss, nil
}

func getClient() http.Client {
	client := http.Client{}
	return client
}

func getFeedRequest(ctx context.Context, url string) (http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return http.Request{}, err
	}
	request.Header.Add("User-Agent", "gator")
	return *request, nil
}

func getFeedResponse(req http.Request) (http.Response, error) {
	client := getClient()
	response, err := client.Do(&req)
	if err != nil {
		return http.Response{}, err
	}
	return *response, nil
}

func readFeedResponse(res http.Response) ([]byte, error) {
	defer res.Body.Close()
	rssData, err := io.ReadAll(res.Body)
	return rssData, err
}

func xmlUnescape(rss *RSSFeed) {
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(rss.Channel.Item[i].Title)
		rss.Channel.Item[i].Description = html.UnescapeString(rss.Channel.Item[i].Description)
	}
}
