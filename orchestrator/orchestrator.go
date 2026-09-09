package orchestrator

import (
	"log"
	deduper "senior_intern_bot/deduper"
	filter "senior_intern_bot/filter"
	poller "senior_intern_bot/poller"
	sender "senior_intern_bot/sender"
	storer "senior_intern_bot/storer"
	"time"
)

type Orchestrator struct {
	poller        poller.Poller
	deduper       deduper.Deduper
	filter        filter.Filterer
	sender        sender.Sender
	storer        storer.Storer
	ticker        *time.Ticker
	isDoneChannel chan bool
}

func New(poller poller.Poller, deduper deduper.Deduper, filter filter.Filterer, sender sender.Sender, storer storer.Storer, ticker *time.Ticker) (Orchestrator, error) {
	return Orchestrator{
		poller:        poller,
		deduper:       deduper,
		filter:        filter,
		sender:        sender,
		storer:        storer,
		ticker:        ticker,
		isDoneChannel: make(chan bool),
	}, nil
}

func (orchestrator *Orchestrator) MainLoop() {
	for {
		select {
		case <-orchestrator.isDoneChannel:
			return
		case <-orchestrator.ticker.C:
			// Poll
			polledPostings, err := orchestrator.poller.Poll()
			if err != nil {
				log.Println("error polling: ", err)
			}

			// Dedupe
			dedupedPostings, err := orchestrator.deduper.Dedupe(polledPostings)
			if err != nil {
				log.Println("error deduping: ", err)
			}

			// Filter
			filteredPostingsByVerdict, err := orchestrator.filter.Filter(dedupedPostings)
			if err != nil {
				log.Println("error filtering: ", err)
			}

			// Sender
			sentPostings, err := orchestrator.sender.Send(filteredPostingsByVerdict)
			if err != nil {
				log.Println("error sending: ", err)
			}

			// Storer (to implement)
			err = orchestrator.storer.Store(sentPostings)
			if err != nil {
				log.Println("error storing: ", err)
			}
		}
	}
}

func (orchestrator *Orchestrator) Start() error {
	go orchestrator.MainLoop()
	return nil
}

func (orchestrator *Orchestrator) Stop() error {
	orchestrator.ticker.Stop()
	orchestrator.isDoneChannel <- true
	return nil
}
