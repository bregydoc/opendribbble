package main

import (
	"fmt"
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

	totalShots := make([]*GenericShot, 0)

	shots := CollectAllPopularShots()
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

	feed, err := GetFeedFromKeyword("machine learning artificial intelligence", map[string]string{
		"max_results": "16",
	})
	if err != nil {
		return totalShots, err
	}

	colors, err := GetPopularColorsFromColorHunt()
	if err != nil {
		return totalShots, err
	}

	if colors == nil {
		colors = make([]*ColorHuntColor, 0)
		colors = append(colors, &ColorHuntColor{
			Color1: "#35013f",
			Color2: "#99ddcc",
			Color3: "#35013f",
			Color4: "#99ddcc",
		})

	} else if len(colors) == 0 {
		colors = append(colors, &ColorHuntColor{
			Color1: "#35013f",
			Color2: "#99ddcc",
			Color3: "#35013f",
			Color4: "#99ddcc",
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
		y, m, d := s.Updated.Date()
		totalShots = append(totalShots, &GenericShot{
			ID:        s.ID,
			Type:      "paper",
			Title:     s.Title,
			Image:     c1 + "," + c2,
			Link:      s.Links[0].Href,
			Comment:   s.Authors[0].Name + fmt.Sprintf(" [%02d/%02d/%02d]", d, m, y),
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
