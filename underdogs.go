package main

import log "github.com/sirupsen/logrus"

func seedUnderdogs(in chan<- scrapeOrder) (err error) {
	log.Info("Underdogs - scrape impossible")
	return
}

// func scrapeUnderdogsEvents() (events []event, err error) {
// 	log.Info("Underdogs - scrape impossible")
// 	return
// }
