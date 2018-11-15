package main

import (
	"encoding/xml"
	"fmt"
	"gopkg.in/resty.v1"
	"log"
	"strings"
	"time"
)

type Feed struct {
	XMLName xml.Name     `xml:"feed"`
	Papers  []*PaperInfo `xml:"entry"`
}

type PaperInfo struct {
	ID        string    `xml:"id"`
	Updated   time.Time `xml:"updated"`
	Published time.Time `xml:"published"`
	Title     string    `xml:"title"`
	Summary   string    `xml:"summary"`
	Authors   []*Author `xml:"author"`
	Links     []*Link   `xml:"link"`
}
type Link struct {
	Href  string `xml:"href,attr"`
	Type  string `xml:"type,attr"`
	Title string `xml:"title,attr"`
	Rel   string `xml:"rel, attr"`
}

type Author struct {
	Name string `xml:"name"`
}

const uriFeed = "https://export.arxiv.org/api/query?search_query=all:"

func GetFeedFromKeyword(keyword string, extraParams ...map[string]string) (*Feed, error) {
	// machine+learning
	keyword = strings.TrimSpace(keyword)
	keyword = strings.Replace(keyword, " ", "+", -1)
	finalUri := uriFeed + keyword
	if len(extraParams) > 0 {
		for name, value := range extraParams[0] {
			finalUri += fmt.Sprintf("&%s=%s", name, value)
		}
	}

	log.Println("Getting from ", finalUri)
	resp, err := resty.R().Get(finalUri)
	if err != nil {
		return nil, err
	}

	data := resp.Body()
	papers := new(Feed)
	err = xml.Unmarshal(data, papers)
	if err != nil {
		return nil, err
	}
	return papers, nil
}
