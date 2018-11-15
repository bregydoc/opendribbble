package main

import (
	"math/rand"
	"strings"
	"time"
)

type GenericShot struct {
	ID        string    `json:"id" storm:"index,unique"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	Link      string    `json:"link" storm:"unique"`
	Comment   string    `json:"comment"`
	Published time.Time `json:"published"`
	Updated   time.Time `xml:"updated"`
}

type ShotsPack struct {
	ID    string         `json:"id" storm:"unique,index"`
	Shots []*GenericShot `json:"shots"`
}

func GetAllShotsFromInternet() ([]*GenericShot, error) {
	rand.Seed(time.Now().Unix())
	shots := CollectAllPopularShots()

	feed, err := GetFeedFromKeyword("machine learning", map[string]string{
		"max_results": "16",
	})

	if err != nil {
		return nil, err
	}

	colors, err := GetPopularColorsFromColorHunt()
	if err != nil {
		return nil, err
	}

	totalShots := make([]*GenericShot, 0)

	for _, s := range shots {
		totalShots = append(totalShots, &GenericShot{
			ID:        s.Link,
			Type:      "art",
			Title:     s.Title,
			Image:     s.Image,
			Link:      s.Link,
			Comment:   s.Comment,
			Published: time.Now(),
			Updated:   time.Now(),
		})
	}

	for _, s := range feed.Papers {
		color := colors[rand.Intn(len(colors))]
		c1, c2, err := GetBestTwoContrastColors(color)
		if err != nil {
			return nil, err
		}
		if rand.Intn(2) == 0 {
			aux := c1
			c1 = c2
			c2 = aux
		}
		totalShots = append(totalShots, &GenericShot{
			ID:        s.ID,
			Type:      "paper",
			Title:     s.Title,
			Image:     c1 + "," + c2,
			Link:      s.Links[0].Href,
			Comment:   s.Authors[0].Name,
			Published: s.Published,
			Updated:   s.Updated,
		})
	}

	return totalShots, nil
}

func FetchAndUpdateShotsOnDB() ([]*GenericShot, error) {
	shots, err := GetAllShotsFromInternet()
	if err != nil {
		return nil, err
	}

	pack := &ShotsPack{
		ID:    "current_shots",
		Shots: shots,
	}

	retriedPack := new(ShotsPack)
	err = ShotsDB.One("ID", "current_shots", retriedPack)

	if err != nil {
		if strings.Contains(err.Error(), "found") {
			err = ShotsDB.Save(pack)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	}

	err = ShotsDB.Update(pack)
	if err != nil {
		return nil, err
	}

	return pack.Shots, nil

}
