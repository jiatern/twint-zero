package Core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	Client *http.Client = new(http.Client)
)

func Request(Query *string, Instance *string, cursor *string) io.ReadCloser {
	var url string = fmt.Sprintf("http://%s/search?f=tweet&q=%s", *Instance, *Query)
	if *cursor != "" {
		url = fmt.Sprintf("http://%s/search%s", *Instance, *cursor)
	}
	hlsCookie := &http.Cookie{
		Name:   "hlsPlayback",
		Value:  "on",
		MaxAge: 300,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.AddCookie(hlsCookie)
	if err != nil {
		log.Fatalf("[nr] %s\n", err)
	}

	req.Header.Set("User-Agent", "twint-zero")
	res, err := Client.Do(req)
	if err != nil {
		log.Fatalf("[do] %s\n", err)
	}

	if res.StatusCode != 200 {
		if 500 <= res.StatusCode && res.StatusCode <= 599 {
			time.Sleep(10 * time.Second)
			return Request(Query, Instance, cursor)
		} else {
			log.Fatalf("status code error: %d %s \n %s", res.StatusCode, res.Status, url)
		}
	}
	return res.Body
}
