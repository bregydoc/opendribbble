package main

import (
	"encoding/json"
	"gopkg.in/go-playground/colors.v1"
	"gopkg.in/resty.v1"
	"math"
	"strings"
)

const baseUri = "https://colorhunt.co/hunt.php"

type ColorHuntColor struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	Likes  string `json:"likes"`
	Code   string `json:"code"`
	Color1 string `json:"color_1"`
	Color2 string `json:"color_2"`
	Color3 string `json:"color_3"`
	Color4 string `json:"color_4"`
}

func Luminance(color *colors.RGBColor) float64 {
	r := float64(color.R) / 255.0
	if r <= 0.03928 {
		r = r / 12.92
	} else {
		r = math.Pow((r+0.055)/1.055, 2.4)
	}

	g := float64(color.G) / 255.0
	if g <= 0.03928 {
		g = g / 12.92
	} else {
		g = math.Pow((g+0.055)/1.055, 2.4)
	}

	b := float64(color.B) / 255.0
	if b <= 0.03928 {
		b = b / 12.92
	} else {
		b = math.Pow((b+0.055)/1.055, 2.4)
	}

	return r*0.2126 + g*0.7152 + b*0.0722
}

func Contrast(color1, color2 *colors.RGBColor) float64 {
	l1 := Luminance(color1) + 0.05
	l2 := Luminance(color2) + 0.05
	return math.Max(l1, l2) / math.Min(l1, l2)
}

func GetBestTwoContrastColors(color *ColorHuntColor) (string, string, error) {
	c1, err := colors.ParseHEX(color.Color1)
	if err != nil {
		return "", "", err
	}
	c1rgb := c1.ToRGB()
	c2, err := colors.ParseHEX(color.Color2)
	if err != nil {
		return "", "", err
	}
	c2rgb := c2.ToRGB()
	c3, err := colors.ParseHEX(color.Color3)
	if err != nil {
		return "", "", err
	}
	c3rgb := c3.ToRGB()
	c4, err := colors.ParseHEX(color.Color4)
	if err != nil {
		return "", "", err
	}
	c4rgb := c4.ToRGB()

	totalColors := []*colors.RGBColor{c1rgb, c2rgb, c3rgb, c4rgb}
	max := 0
	contrasts := []float64{0.0, 0.0, 0.0}
	for i, c := range totalColors[1:] {
		contrasts[i] = Contrast(totalColors[0], c)
		if contrasts[max] < contrasts[i] {
			max = i
		}
	}

	return totalColors[0].ToHEX().String(), totalColors[max+1].ToHEX().String(), nil

}

func GetPopularColorsFromColorHunt() ([]*ColorHuntColor, error) {
	resp, err := resty.R().SetFormData(map[string]string{
		"sort": "popular",
	}).Post(baseUri)
	if err != nil {
		return nil, err
	}

	data := string(resp.Body())
	data = strings.Replace(data, "<script>arr = ", `{"data":`, -1)
	data = strings.Replace(data, ", ];</script>", `]}`, -1)

	//log.Println(data)
	//colors := make([]*ColorHuntColor, 0)

	type Response struct {
		Data []*ColorHuntColor `json:"data"`
	}
	tColors := Response{}
	err = json.Unmarshal([]byte(data), &tColors)
	if err != nil {
		return nil, err
	}

	for _, c := range tColors.Data {
		c.Color1 = "#" + c.Code[:6]
		c.Color2 = "#" + c.Code[6:12]
		c.Color3 = "#" + c.Code[12:18]
		c.Color4 = "#" + c.Code[18:]
	}

	return tColors.Data, nil
}
