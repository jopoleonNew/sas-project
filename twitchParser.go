package main

import (
	"log"

	"net/http"

	json2 "encoding/json"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	//selector
	//#ember147553 > div > div.flex.flex-grow-1.flex-nowrap.align-items-center.mg-r-1 > div > h3
	//#ember1906 > div > div.flex.flex-grow-1.flex-nowrap.align-items-center.mg-r-1 > div > h3 class="font-size-4 mg-b-05 js-card__title qa-card__title"
	apiURL := "https://api.twitch.tv/kraken/streams/?game=Dota%202"
	client := &http.Client{}
	r, _ := http.NewRequest("GET", apiURL, nil)
	r.Header.Add("Client-ID", "j15uy1j9edpubl48bta3733qepj8ck")

	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)

	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)

	}
	//log.Println(string(body))
	var t TwitchResponse
	err = json2.Unmarshal(body, &t)
	if err != nil {
		log.Fatal(err)

	}
	for _, stream := range t.Streams {
		log.Println(stream.Channel.Name)
		log.Println(stream.StreamType)
		log.Println(stream.Channel.Status)
		name := stream.Channel.Status
		if strings.Contains(name, "Virtus Pro") || strings.Contains(name, "Virtus") || strings.Contains(name, "virtus") {
			//some logic you want
		}
	}
}

type TwitchResponse struct {
	Total   int `json:"_total"`
	Streams []struct {
		ID          int64     `json:"_id"`
		Game        string    `json:"game"`
		Viewers     int       `json:"viewers"`
		VideoHeight int       `json:"video_height"`
		AverageFps  float64   `json:"average_fps"`
		Delay       int       `json:"delay"`
		CreatedAt   time.Time `json:"created_at"`
		IsPlaylist  bool      `json:"is_playlist"`
		StreamType  string    `json:"stream_type"`
		Preview     struct {
			Small    string `json:"small"`
			Medium   string `json:"medium"`
			Large    string `json:"large"`
			Template string `json:"template"`
		} `json:"preview"`
		Channel struct {
			Mature                       bool        `json:"mature"`
			Partner                      bool        `json:"partner"`
			Status                       string      `json:"status"`
			BroadcasterLanguage          string      `json:"broadcaster_language"`
			DisplayName                  string      `json:"display_name"`
			Game                         string      `json:"game"`
			Language                     string      `json:"language"`
			ID                           int         `json:"_id"`
			Name                         string      `json:"name"`
			CreatedAt                    time.Time   `json:"created_at"`
			UpdatedAt                    time.Time   `json:"updated_at"`
			Delay                        interface{} `json:"delay"`
			Logo                         string      `json:"logo"`
			Banner                       interface{} `json:"banner"`
			VideoBanner                  string      `json:"video_banner"`
			Background                   interface{} `json:"background"`
			ProfileBanner                string      `json:"profile_banner"`
			ProfileBannerBackgroundColor string      `json:"profile_banner_background_color"`
			URL                          string      `json:"url"`
			Views                        int         `json:"views"`
			Followers                    int         `json:"followers"`
			Links                        struct {
				Self          string `json:"self"`
				Follows       string `json:"follows"`
				Commercial    string `json:"commercial"`
				StreamKey     string `json:"stream_key"`
				Chat          string `json:"chat"`
				Features      string `json:"features"`
				Subscriptions string `json:"subscriptions"`
				Editors       string `json:"editors"`
				Teams         string `json:"teams"`
				Videos        string `json:"videos"`
			} `json:"_links"`
		} `json:"channel"`
		Links struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"streams"`
	Links struct {
		Self     string `json:"self"`
		Next     string `json:"next"`
		Featured string `json:"featured"`
		Summary  string `json:"summary"`
		Followed string `json:"followed"`
	} `json:"_links"`
}
