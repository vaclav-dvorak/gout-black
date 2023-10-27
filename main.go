package main

import (
	"os"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
)

const (
	workers = 5
)

var (
	wg sync.WaitGroup
)

type event struct {
	title string
	date  time.Time
	desc  string
	score int
	venue string
}

type scrapeOrder struct {
	url      string
	scrapper func(url string, output chan<- event) error
	title    string
}

func scrape(in <-chan scrapeOrder, out chan<- event) {
	defer wg.Done()
	for input := range in {
		start := time.Now()
		log.Info("Request - " + input.title)
		if err := input.scrapper(input.url, out); err != nil {
			log.Fatal(err)
		}
		log.Infof("Request - %s - done - took %s", input.title, time.Since(start))
	}
}

func main() {
	bufScrapeOrder := make(chan scrapeOrder, workers)
	bufEvent := make(chan event, workers)
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go scrape(bufScrapeOrder, bufEvent)
	}

	events := make([]event, 0)
	go func() {
		for event := range bufEvent {
			events = append(events, event)
		}
	}()

	seeders := []func(chan<- scrapeOrder) error{seedVopice, seedFuturum, seedUnderdogs, seedAkropolis}
	for _, seeder := range seeders {
		if err := seeder(bufScrapeOrder); err != nil {
			log.Fatal(err)
		}
	}
	close(bufScrapeOrder)

	wg.Wait()
	close(bufEvent)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"date", "venue", "title", "desc", "score"})
	for _, v := range events {
		t.AppendRow(table.Row{
			v.date.Format("02.01.2006"), v.venue, v.title, v.desc, v.score,
		})
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 4, WidthMax: 90, WidthMaxEnforcer: text.WrapSoft},
	})
	t.SortBy([]table.SortBy{
		{Name: "date", Mode: table.Asc},
		{Name: "score", Mode: table.Dsc},
	})
	t.SetStyle(table.StyleLight)
	t.Render()
}
