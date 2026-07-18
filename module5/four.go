package module5

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func MultiURLTime(urls []string) {
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			URLTime(u)
		}(url)
	}

	wg.Wait()
}

func URLTime(url string) {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error: %q - %s", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("error: %q - bad status - %s", url, resp.Status)
		return
	}

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		log.Printf("error: %q - %s", url, err)
		return
	}

	time.Sleep(2 * time.Second)

	duration := time.Since(start)
	log.Printf("info: %q - %v", url, duration)
}
