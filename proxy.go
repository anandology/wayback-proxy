package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type WaybackResponse struct {
	URL               string                   `json:"url"`
	Timestamp         string                   `json:"timestamp"`
	ArchivedSnapshots WaybackResponseSnapshots `json:"archived_snapshots"`
}

type WaybackResponseSnapshots struct {
	Closest *WaybackResponseEntry `json:"closest,omitempty"`
}

type WaybackResponseEntry struct {
	Timestamp string `json:"timestamp"`
	Available bool   `json:"available"`
	Status    string `json:"status"`
	URL       string `json:"url"`
}

func GetWaybackURL(requestURL string, timestamp string) (string, error) {
	apiURL := url.URL{
		Scheme: "https",
		Host:   "archive.org",
		Path:   "/wayback/available",
	}

	q := url.Values{}
	q.Set("timestamp", timestamp)
	q.Set("url", requestURL)
	apiURL.RawQuery = q.Encode()

	log.Println("API URL", apiURL.String())

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var out WaybackResponse
	err = json.Unmarshal(body, &out)
	if err != nil {
		return "", err
	}

	if out.ArchivedSnapshots.Closest == nil {
		return "", nil
	}

	log.Println("JSON", out.ArchivedSnapshots.Closest)

	wbURL := fmt.Sprintf(
		"https://web.archive.org/web/%sid_/%s",
		out.ArchivedSnapshots.Closest.Timestamp,
		out.URL)
	return wbURL, nil
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("GET", req.URL.String())
		url, err := GetWaybackURL(req.URL.String(), "2015")

		if err != nil {
			internalError(w, err)
			return
		}
		if url == "" {
			w.WriteHeader(404)
			w.Write([]byte("<h1>Not Found</b>"))
		}

		resp, err := http.Get(url)
		if err != nil {
			internalError(w, err)
			return
		}

		read, write := io.Pipe()
		go func() {
			defer write.Close()
			defer resp.Body.Close()
			io.Copy(write, resp.Body)
		}()
		io.Copy(w, read)
	})

	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func internalError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte("<h1>Internal Error</b>"))
	w.Write([]byte(err.Error()))
}
