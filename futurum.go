package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func seedFuturum(in chan<- scrapeOrder) (err error) {
	req, terr := http.NewRequest(http.MethodGet, "https://futurum.musicbar.cz/program/", nil)
	if terr != nil {
		err = terr
		return
	}
	req.Header.Set("User-Agent", "ScraperBot - We read events once a day")

	res, terr := http.DefaultClient.Do(req)
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

	candidates := make([]string, 0)
	doc.Find(".event-lines").Each(func(i int, s *goquery.Selection) {
		cand, _ := s.Attr("href")
		candidates = append(candidates, cand)
	})

	log.Infof("Futurum - found %d events", len(candidates))
	for k, v := range candidates {
		in <- scrapeOrder{url: v, scrapper: scrapeFuturumEvents, title: fmt.Sprintf("Futurum [%d/%d]", k+1, len(candidates))}
	}
	return
}

func scrapeFuturumEvents(url string, eventChan chan<- event) (err error) {
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

	title := doc.Find(".single-blok h1").Text()
	date := doc.Find(".single-blok .block span").Eq(0).Text() + " @ " + doc.Find(".single-blok .block span").Eq(2).Text()
	dateParsed, err := time.Parse("2.1.2006 @ 15:04", date)
	if err != nil {
		log.Infof("Cannot parse datetime: %s\n", date)
	}
	desc := doc.Find(".event_content").Text()
	desc = strings.Replace(strings.TrimSpace(desc), "\n\n", "\n", -1) //? sanitize description by removing trailing spaces and empty lines
	score := strings.Count(desc, "black")
	if score > 1 {
		eventChan <- event{title: title, date: dateParsed, desc: desc, score: score, venue: "Futurum musicbar"}
	}

	return
}
