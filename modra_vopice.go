package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func seedVopice(in chan<- scrapeOrder) (err error) {
	const pageCount = 4
	for i := 1; i <= pageCount; i++ {
		in <- scrapeOrder{url: fmt.Sprintf("https://modravopice.eu/program/nadchazejici/?action=tribe_list&tribe_paged=%d", i), scrapper: scrapeVopiceEvents, title: fmt.Sprintf("Modra vopice [%d/4]", i)}
	}
	return
}

func scrapeVopiceEvents(url string, eventChan chan<- event) (err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "ScraperBot - We read events once a day")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	doc.Find("#tribe-events-content .type-tribe_events").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".tribe-events-list-event-title a").Text())
		date := s.Find(".date-start").Text()
		dateParsed, err := time.Parse("2.1.2006 @ 15:04", date)
		if err != nil {
			log.Infof("Cannot parse datetime: %s\n", date)
		}
		desc := s.Find(".tribe-events-list-event-description p").Text()
		desc = strings.Replace(strings.TrimSpace(desc), "\n", "\n\n", -1) //? sanitize description by removing trailing spaces and empty lines
		score := strings.Count(desc, "black")
		if score > 1 {
			eventChan <- event{title: title, date: dateParsed, desc: desc, score: score, venue: "Modrá vopice"}
		}
	})
	return
}
