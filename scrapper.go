package main

import (
	"strings"

	"github.com/gocolly/colly"
)

// DribbbleShot represents the dribbble shot structure
type DribbbleShot struct {
	Title   string `json:"title"`
	Image   string `json:"image"`
	Link    string `json:"link"`
	Comment string `json:"comment"`

	Extras ExtraInfo `json:"extras"`
}

// ExtraInfo represents the extra information of the shot
type ExtraInfo struct {
	Team  string `json:"team"`
	User  string `json:"user"`
	Fav   string `json:"fav"`
	Cmnt  string `json:"cmnt"`
	Views string `json:"views"`
}

// CollectAllPopularShots ...
func CollectAllPopularShots() []*DribbbleShot {
	shots := []*DribbbleShot{}
	c := colly.NewCollector()

	c.OnHTML("ol.dribbbles", func(e *colly.HTMLElement) {
		validItems := []*colly.HTMLElement{}
		e.ForEach("li", func(i int, item *colly.HTMLElement) {
			if strings.Contains(item.Attr("id"), "screenshot") {

				validItems = append(validItems, item)
			}
		})

		for _, item := range validItems {

			over := item.DOM.Find("a.dribbble-over").First()
			title := over.Find("strong").Text()
			comment := over.Find("span.comment").Text()

			prel := item.DOM.Find("a.dribbble-link").First()

			image, _ := prel.Find("source").Attr("srcset")
			link, _ := prel.Attr("href")

			teamChunk := item.DOM.Find("span.attribution-team").First()
			team := teamChunk.Find("a").First().Text()

			userChunk := item.DOM.Find("span.attribution-user").First()
			user := userChunk.Find("a").First().Text()

			tools := item.DOM.Find("ul.tools").First()
			fav := tools.Find("li.fav").First().Text()
			cmnt := tools.Find("li.cmnt").First().Text()
			views := tools.Find("li.views").First().Text()

			team = strings.TrimSpace(team)
			user = strings.TrimSpace(user)

			fav = strings.Replace(fav, "\n", "", -1)
			fav = strings.Replace(fav, ",", "", -1)
			fav = strings.TrimSpace(fav)
			cmnt = strings.Replace(cmnt, "\n", "", -1)
			cmnt = strings.Replace(cmnt, ",", "", -1)
			cmnt = strings.TrimSpace(cmnt)
			views = strings.Replace(views, "\n", "", -1)
			views = strings.Replace(views, ",", "", -1)
			views = strings.TrimSpace(views)

			shots = append(shots, &DribbbleShot{
				Title:   title,
				Comment: comment,
				Image:   image,
				Link:    link,
				Extras: ExtraInfo{
					Team:  team,
					User:  user,
					Fav:   fav,
					Cmnt:  cmnt,
					Views: views,
				},
			})

		}
	})

	c.Visit(DribbleShots)

	return shots
}
