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

func scrapeVopiceEvents(url string, eventChan chan<- event) (events []event, err error) {
	req, terr := http.NewRequest(http.MethodGet, url, nil)
	if terr != nil {
		err = terr
		return
	}
	req.Header.Set("User-Agent", "ScraperBot - We read events once a day")

	client := &http.Client{}
	res, err := client.Do(req)
	if terr != nil {
		err = terr
		return
	}
	defer res.Body.Close()
	doc, terr := goquery.NewDocumentFromReader(res.Body)
	if terr != nil {
		err = terr
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
		score := strings.Count(desc, "black")
		if score > 1 {
			eventChan <- event{title: title, date: dateParsed, desc: desc, score: score, venue: "Modrá vopice"}
		}
	})
	return
}
